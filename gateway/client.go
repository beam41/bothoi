package gateway

import (
	"bothoi/models"
	"bothoi/references/gateway_opcode"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// connect to the discord gateway.
// this must called last in main because infinite loop that keep heartbeat.
func Connect() {
	c, _, err := websocket.DefaultDialer.Dial(os.Getenv("GATEWAY_URL"), nil)
	if err != nil {
		log.Panicln(err)
	}

	c.WriteJSON(models.NewIdentify())

	heatbeatInterval := make(chan int)

	// receive the gateway response
	go func() {
		for {
			var payload models.GatewayPayload
			err := c.ReadJSON(&payload)
			if err != nil {
				log.Panicln(err)
				return
			}
			log.Println("incoming: ", payload)
			switch payload.Op {
			case gateway_opcode.Hello:
				heatbeatInterval <- int(payload.D.(map[string]interface{})["heartbeat_interval"].(float64))
			}
			setSequenceNumber(payload.S)
		}
	}()

	// keeping heartbeats and prevent application from closing.
	interval := <-heatbeatInterval

	time.Sleep(time.Duration(float64(interval)*rand.Float64()) * time.Millisecond)
	WriteJSONLog(c, models.NewHeartbeat(0))

	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		sequenceNumberLock.Lock()
		WriteJSONLog(c, models.NewHeartbeat(*sequenceNumber))
		sequenceNumberLock.Unlock()
	}
}

func WriteJSONLog(c *websocket.Conn, v interface{}) error {
	log.Println("outgoing: ", v)
	return c.WriteJSON(v)
}
