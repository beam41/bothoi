package gateway

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/references/gateway_opcode"
	"bothoi/states"
	"bothoi/util/ws_util"
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// connect to the discord gateway.
func Connect() {
	resume := false
	for {
		connection(resume)
		resume = true
	}
}

func connection(isResume bool) {
	c, _, err := websocket.DefaultDialer.Dial(config.GATEWAY_URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	if !isResume {
		ws_util.WriteJSONLog(c, models.NewIdentify())
		states.SessionStateReady.Add(1)
	} else {
		ws_util.WriteJSONLog(c, models.NewResume(sequenceNumber, states.SessionState.SessionID))
	}

	heatbeatInterval := make(chan int)
	heatbeatAcked := make(chan bool)
	immediateHeartbeat := make(chan bool)

	// receive the gateway response
	go func() {
		for {
			var payload models.GatewayPayload
			err := c.ReadJSON(&payload)
			if err != nil {
				log.Println(err)
				if strings.HasPrefix(err.Error(), "websocket: close 1001") {
					ws_util.WriteJSONLog(c, models.NewIdentify())
					states.SessionStateReady.Add(1)
				}
				continue
			}
			if (config.DEVELOPMENT) {
				jsonDat, _ := json.Marshal(payload)
				log.Println("incoming: ", payload, string(jsonDat))
			} else {
				log.Println("incoming: ", payload)
			}

			// log.Println("incoming: ", payload)
			setSequenceNumber(payload.S)
			switch payload.Op {
			case gateway_opcode.Hello:
				heatbeatInterval <- int(payload.D.(map[string]interface{})["heartbeat_interval"].(float64))
			case gateway_opcode.HeartbeatAck:
				heatbeatAcked <- true
			case gateway_opcode.Heartbeat:
				immediateHeartbeat <- true
			case gateway_opcode.Dispatch:
				go dispatchHandler(c, payload)
			case gateway_opcode.InvalidSession:
				ws_util.WriteJSONLog(c, models.NewIdentify())
				states.SessionStateReady.Add(1)
			}
		}
	}()

	// keeping heartbeats and prevent application from closing.
	interval := <-heatbeatInterval

	time.Sleep(time.Duration(float64(interval)*rand.Float64()) * time.Millisecond)
	ws_util.WriteJSONLog(c, models.NewHeartbeat(nil))

	for {
		// wait for heartbeat ack
		select {
		case <-heatbeatAcked:
		case <-time.After(time.Duration(interval) * time.Millisecond):
			// uh oh timeout, reconnect
			log.Println("timeout, attempting to reconnect")
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseServiceRestart, ""))
			return
		}
		// wait for next heartbeat
		select {
		case <-immediateHeartbeat:
		case <-time.After(time.Duration(interval) * time.Millisecond):
		}
		sequenceNumberLock.Lock()
		ws_util.WriteJSONLog(c, models.NewHeartbeat(sequenceNumber))
		sequenceNumberLock.Unlock()
	}
}


