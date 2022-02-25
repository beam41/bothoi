package models

import (
	"bothoi/references/gateway_opcode"
	"log"
	"os"
	"strconv"
)

type GatewayPayload struct {
	Op uint8       `json:"op"`
	S  *uint64     `json:"s"`
	T  string      `json:"t"`
	D  interface{} `json:"d"`
}

func NewHeartbeat(s uint64) GatewayPayload {
	var d interface{}
	if s != 0 {
		d = s
	}

	return GatewayPayload{
		Op: gateway_opcode.Heartbeat,
		D:  d,
	}
}

func NewIdentify() GatewayPayload {
	intent, err := strconv.Atoi(os.Getenv("GATEWAY_INTENT"))
	if err != nil {
		log.Fatalln("Intent is not a number")
	}
	return GatewayPayload{
		Op: gateway_opcode.Identify,
		D: map[string]interface{}{
			"token":   os.Getenv("TOKEN"),
			"intents": intent,
			"properties": map[string]interface{}{
				"$os":      "linux",
				"$browser": "bothoi",
				"$device":  "bothoi",
			},
			"compress": false,
			"shard":    []int{0, 1},
		},
	}
}
