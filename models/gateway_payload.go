package models

import (
	"bothoi/references/gateway_opcode"
	"log"
	"os"
	"strconv"
)

type GatewayPayload struct {
	Op uint8       `json:"op"`
	S  *uint64     `json:"s,omitempty"`
	T  string      `json:"t"`
	D  interface{} `json:"d"`
}

func NewHeartbeat(s *uint64) GatewayPayload {
	return GatewayPayload{
		Op: gateway_opcode.Heartbeat,
		D:  s,
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
			"token":   os.Getenv("BOT_TOKEN"),
			"intents": intent,
			"properties": map[string]interface{}{
				"$os":      "linux",
				"$browser": "bothoi",
				"$device":  "bothoi",
			},
			"compress": false,
		},
	}
}

func NewResume(s *uint64, sessionId string) GatewayPayload {
	return GatewayPayload{
		Op: gateway_opcode.Resume,
		D: map[string]interface{}{
			"token":      os.Getenv("BOT_TOKEN"),
			"session_id": sessionId,
			"seq":        s,
		},
	}
}
