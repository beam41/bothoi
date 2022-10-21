package global

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var conn *websocket.Conn

func StartGatewayConn() {
	c, _, err := websocket.DefaultDialer.Dial(config.GatewayUrl, nil)
	conn = c
	if err != nil {
		log.Fatalln(err)
	}
}

func CloseGatewayConn() {
	err := conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func GatewayConnReadJSON(v any) (err error) {
	err = conn.ReadJSON(&v)
	return err
}

func GatewayConnWriteJSON(v any) (err error) {
	log.Println("outgoing: ", v)
	err = conn.WriteJSON(v)
	return
}

func GatewayConnCloseRestart() {
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseServiceRestart, ""))
	if err != nil {
		log.Println(err)
	}
}

var sequenceNumber struct {
	sync.RWMutex
	n *uint64
}

func GetSequenceNumber() *uint64 {
	return sequenceNumber.n
}

func SequenceNumberRLock() {
	sequenceNumber.RLock()
}

func SequenceNumberRUnLock() {
	sequenceNumber.RUnlock()
}

func SetSequenceNumber(n *uint64) {
	sequenceNumber.Lock()
	sequenceNumber.n = n
	sequenceNumber.Unlock()
}

var gatewaySession struct {
	sync.RWMutex
	state *discord_models.ReadyEvent
}

func AddGatewaySession(state *discord_models.ReadyEvent) {
	gatewaySession.Lock()
	gatewaySession.state = state
	gatewaySession.Unlock()
}

func GetGatewaySession() *discord_models.ReadyEvent {
	gatewaySession.RLock()
	defer gatewaySession.RUnlock()
	return gatewaySession.state
}
