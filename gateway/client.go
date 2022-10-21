package gateway

import (
	"bothoi/config"
	"bothoi/global"
	"bothoi/models"
	"bothoi/references/gateway_opcode"
	"bothoi/repo"
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"
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
	global.StartGatewayConn()
	defer global.CloseGatewayConn()

	if !isResume {
		err := global.GatewayConnWriteJSON(models.NewIdentify())
		if err != nil {
			log.Println(err)
		}
	} else {
		err := global.GatewayConnWriteJSON(models.NewResume(global.GetSequenceNumber(), repo.GetSessionState().SessionID))
		if err != nil {
			log.Println(err)
		}
	}

	heatbeatInterval := make(chan int)
	heatbeatAcked := make(chan struct{})
	immediateHeartbeat := make(chan struct{})

	// receive the gateway response
	go func() {
		for {
			var payload models.GatewayPayload
			err := global.GatewayConnReadJSON(&payload)
			if err != nil {
				log.Println(err)
				if strings.HasPrefix(err.Error(), "websocket: close 1001") {
					err := global.GatewayConnWriteJSON(models.NewIdentify())
					if err != nil {
						log.Println(err)
						return
					}
				}
				continue
			}
			if config.Development {
				jsonDat, _ := json.Marshal(payload)
				log.Println("incoming: ", payload, string(jsonDat))
			} else {
				log.Println("incoming: ", payload)
			}

			// log.Println("incoming: ", payload)
			if payload.S != nil {
				global.SetSequenceNumber(payload.S)
			}
			switch payload.Op {
			case gateway_opcode.Hello:
				heatbeatInterval <- int(payload.D.(map[string]any)["heartbeat_interval"].(float64))
			case gateway_opcode.HeartbeatAck:
				heatbeatAcked <- struct{}{}
			case gateway_opcode.Heartbeat:
				immediateHeartbeat <- struct{}{}
			case gateway_opcode.Dispatch:
				go dispatchHandler(payload)
			case gateway_opcode.Reconnect:
				fallthrough
			case gateway_opcode.InvalidSession:
				global.GatewayConnCloseRestart()
				return
			}
		}
	}()

	// keeping heartbeats and prevent application from closing.
	interval := <-heatbeatInterval

	time.Sleep(time.Duration(float64(interval)*rand.Float64()) * time.Millisecond)
	err := global.GatewayConnWriteJSON(models.NewHeartbeat(nil))
	if err != nil {
		log.Println(err)
	}
	for {
		// wait for heartbeat ack
		select {
		case <-heatbeatAcked:
		case <-time.After(time.Duration(interval) * time.Millisecond):
			// uh oh timeout, reconnect
			log.Println("timeout, attempting to reconnect")
			global.GatewayConnCloseRestart()
			return
		}
		// wait for next heartbeat
		select {
		case <-immediateHeartbeat:
		case <-time.After(time.Duration(interval) * time.Millisecond):
		}
		global.SequenceNumberRLock()
		err := global.GatewayConnWriteJSON(models.NewHeartbeat(global.GetSequenceNumber()))
		if err != nil {
			log.Println(err)
		}
		global.SequenceNumberRUnLock()
	}
}
