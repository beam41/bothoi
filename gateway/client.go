package gateway

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/references/gateway_opcode"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"sync"
	"time"
)

type voiceInstantiateChan struct {
	sessionIDChan   chan<- string
	voiceServerChan chan<- *discord_models.VoiceServer
}

type Client struct {
	conn      *websocket.Conn
	ctx       context.Context
	ctxCancel context.CancelFunc
	info      struct {
		sync.RWMutex
		sequenceNumber *uint64
		session        *discord_models.ReadyEvent
	}
	voiceInstantiateList struct {
		sync.RWMutex
		list map[types.Snowflake]voiceInstantiateChan
	}
	resume                  bool
	interactionExecutorList map[string]func(*discord_models.Interaction)
	newSessionIDHandler     func(types.Snowflake, string)
	waitResume              chan struct{}
}

func NewClient() *Client {
	return &Client{
		voiceInstantiateList: struct {
			sync.RWMutex
			list map[types.Snowflake]voiceInstantiateChan
		}{list: make(map[types.Snowflake]voiceInstantiateChan)},
		waitResume: make(chan struct{}, 1),
	}
}

func (client *Client) gatewayConnReadJSON(v any) (err error) {
	err = client.conn.ReadJSON(&v)
	return err
}

func (client *Client) gatewayConnWriteJSON(v any) (err error) {
	log.Println("outgoing: ", v)
	err = client.conn.WriteJSON(v)
	return
}

func (client *Client) gatewayConnCloseRestart() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("gatewayConnCloseRestart panic occurred:", err)
		}
	}()
	client.ctxCancel()
	err := client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseServiceRestart, ""))
	if err != nil {
		log.Println(err)
	}
	client.resume = true
	client.waitResume <- struct{}{}
}

// connect to the discord gateway.
func (client *Client) Connect() {
	client.waitResume <- struct{}{}
	for {
		<-client.waitResume
		ctx, cancel := context.WithCancel(context.Background())
		client.ctx = ctx
		client.ctxCancel = cancel
		client.connection()
	}
}

func (client *Client) connection() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("gateway connection panic occurred:", err)
		}
	}()
	c, _, err := websocket.DefaultDialer.Dial(config.GatewayUrl, nil)
	client.conn = c
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		err := client.conn.Close()
		if err != nil {
			log.Panicln(err)
		}
	}()

	if !client.resume {
		err := client.gatewayConnWriteJSON(discord_models.NewIdentify())
		if err != nil {
			log.Println(err)
		}
	} else {
		client.info.RLock()
		v := discord_models.NewResume(client.info.sequenceNumber, client.info.session.SessionID)
		log.Println("outgoing: ", v)
		err := client.conn.WriteJSON(v)
		client.info.RUnlock()
		if err != nil {
			log.Println(err)
		}
	}

	heartbeatInterval := make(chan int)
	heartbeatAcked := make(chan struct{})
	immediateHeartbeat := make(chan struct{})

	// receive the gateway response
	go func() {
		for {
			var payload discord_models.GatewayPayload
			err := client.gatewayConnReadJSON(&payload)
			if err != nil {
				log.Println(err)
				if websocket.IsCloseError(err, 1001, 4004, 4010, 4011, 4012, 4013, 4014) {
					client.gatewayConnCloseRestart()
				} else {
					client.ctxCancel()
				}
				return
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
				heartbeatInterval <- int(payload.D.(map[string]any)["heartbeat_interval"].(float64))
			case gateway_opcode.HeartbeatAck:
				heartbeatAcked <- struct{}{}
			case gateway_opcode.Heartbeat:
				immediateHeartbeat <- struct{}{}
			case gateway_opcode.Dispatch:
				go client.dispatchHandler(payload)
			case gateway_opcode.Reconnect:
				client.gatewayConnCloseRestart()
				return
			case gateway_opcode.InvalidSession:
				if payload.D.(bool) {
					client.gatewayConnCloseRestart()
				}
				return
			}
		}
	}()

	// keeping heartbeats and prevent application from closing.
	interval := <-heartbeatInterval

	time.Sleep(time.Duration(float64(interval)*rand.Float64()) * time.Millisecond)
	err = client.gatewayConnWriteJSON(discord_models.NewHeartbeat(nil))
	if err != nil {
		log.Println(err)
	}
	heartbeatIntervalTicker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer heartbeatIntervalTicker.Stop()
	for {
		// wait for heartbeat ack
		select {
		case <-client.ctx.Done():
			return
		case <-heartbeatAcked:
		case <-heartbeatIntervalTicker.C:
			// uh oh timeout, reconnect
			log.Println("timeout, attempting to reconnect")
			client.gatewayConnCloseRestart()
			return
		}

		// wait for next heartbeat
		select {
		case <-client.ctx.Done():
			return
		case <-immediateHeartbeat:
		case <-heartbeatIntervalTicker.C:
		}
		client.info.RLock()
		err := client.gatewayConnWriteJSON(discord_models.NewHeartbeat(client.info.sequenceNumber))
		client.info.RUnlock()
		if err != nil {
			log.Println(err)
		}
	}
}
