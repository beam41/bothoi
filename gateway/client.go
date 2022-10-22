package gateway

import (
	"bothoi/config"
	"bothoi/gateway/gateway_interface"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/references/gateway_opcode"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type voiceInstantiateChan struct {
	sessionIdChan   chan<- string
	voiceServerChan chan<- *discord_models.VoiceServer
}

type client struct {
	conn *websocket.Conn
	info struct {
		sync.RWMutex
		sequenceNumber *uint64
		session        *discord_models.ReadyEvent
	}
	voiceWaiter struct {
		sync.RWMutex
		list map[types.Snowflake]voiceInstantiateChan
	}
}

func NewClient() gateway_interface.ClientInterface {
	return &client{
		voiceWaiter: struct {
			sync.RWMutex
			list map[types.Snowflake]voiceInstantiateChan
		}{list: make(map[types.Snowflake]voiceInstantiateChan)},
	}
}

func (client *client) gatewayConnReadJSON(v any) (err error) {
	err = client.conn.ReadJSON(&v)
	return err
}

func (client *client) gatewayConnWriteJSON(v any) (err error) {
	log.Println("outgoing: ", v)
	err = client.conn.WriteJSON(v)
	return
}

func (client *client) gatewayConnCloseRestart() {
	err := client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseServiceRestart, ""))
	if err != nil {
		log.Println(err)
	}
}

// connect to the discord gateway.
func (client *client) Connect() {
	resume := false
	for {
		client.connection(resume)
		resume = true
	}
}

func (client *client) connection(isResume bool) {
	c, _, err := websocket.DefaultDialer.Dial(config.GatewayUrl, nil)
	client.conn = c
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := client.conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	if !isResume {
		err := client.gatewayConnWriteJSON(discord_models.NewIdentify())
		if err != nil {
			log.Println(err)
		}
	} else {
		client.info.RLock()
		err := client.gatewayConnWriteJSON(discord_models.NewResume(client.info.sequenceNumber, client.info.session.SessionId))
		client.info.RUnlock()
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
			var payload discord_models.GatewayPayload
			err := client.gatewayConnReadJSON(&payload)
			if err != nil {
				log.Println(err)
				if strings.HasPrefix(err.Error(), "websocket: close 1001") {
					err := client.gatewayConnWriteJSON(discord_models.NewIdentify())
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
				client.info.Lock()
				client.info.sequenceNumber = payload.S
				client.info.Unlock()
			}

			switch payload.Op {
			case gateway_opcode.Hello:
				heatbeatInterval <- int(payload.D.(map[string]any)["heartbeat_interval"].(float64))
			case gateway_opcode.HeartbeatAck:
				heatbeatAcked <- struct{}{}
			case gateway_opcode.Heartbeat:
				immediateHeartbeat <- struct{}{}
			case gateway_opcode.Dispatch:
				go client.dispatchHandler(payload)
			case gateway_opcode.Reconnect:
				fallthrough
			case gateway_opcode.InvalidSession:
				client.gatewayConnCloseRestart()
				return
			}
		}
	}()

	// keeping heartbeats and prevent application from closing.
	interval := <-heatbeatInterval

	time.Sleep(time.Duration(float64(interval)*rand.Float64()) * time.Millisecond)
	err = client.gatewayConnWriteJSON(discord_models.NewHeartbeat(nil))
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
			client.gatewayConnCloseRestart()
			return
		}
		// wait for next heartbeat
		select {
		case <-immediateHeartbeat:
		case <-time.After(time.Duration(interval) * time.Millisecond):
		}
		client.info.RLock()
		err := client.gatewayConnWriteJSON(discord_models.NewHeartbeat(client.info.sequenceNumber))
		client.info.RUnlock()
		if err != nil {
			log.Println(err)
		}
	}
}
