package ws_util

import (
	"log"

	"github.com/gorilla/websocket"
)

func WriteJSONLog(c *websocket.Conn, v interface{}) error {
	log.Println("outgoing: ", v)
	return c.WriteJSON(v)
}
