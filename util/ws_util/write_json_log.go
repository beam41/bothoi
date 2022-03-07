package ws_util

import (
	"log"

	"github.com/gorilla/websocket"
)

func WriteJSONLog(c *websocket.Conn, v interface{}, voice bool) (err error) {
	if voice {
		log.Println("outgoing voice: ", v)
	} else {
		log.Println("outgoing: ", v)
	}
	err = c.WriteJSON(v)
	if err != nil {
		log.Println(err)
	}
	return
}
