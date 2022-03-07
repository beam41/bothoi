package models

import (
	"bothoi/config"
	"bothoi/references/gateway_opcode"
	"bothoi/references/voice_opcode"
	"strconv"
)

type GatewayPayload struct {
	Op uint8       `json:"op" mapstructure:"op"`
	S  *uint64     `json:"s,omitempty" mapstructure:"s"`
	T  string      `json:"t" mapstructure:"t"`
	D  interface{} `json:"d" mapstructure:"d"`
}

func NewHeartbeat(s *uint64) GatewayPayload {
	return GatewayPayload{
		Op: gateway_opcode.Heartbeat,
		D:  s,
	}
}

func NewIdentify() GatewayPayload {
	return GatewayPayload{
		Op: gateway_opcode.Identify,
		D: map[string]interface{}{
			"token":   config.BOT_TOKEN,
			"intents": config.GATEWAY_INTENT,
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
			"token":      config.BOT_TOKEN,
			"session_id": sessionId,
			"seq":        s,
		},
	}
}

func NewVoiceStateUpdate(guildId, voiceId string, mute, deaf bool) GatewayPayload {
	return GatewayPayload{
		Op: gateway_opcode.VoiceStateUpdate,
		D: map[string]interface{}{
			"guild_id":   guildId,
			"channel_id": voiceId,
			"self_mute":  mute,
			"self_deaf":  deaf,
		},
	}
}

func NewVoiceIdentify(guildId, userId, sessionId, token string) GatewayPayload {
	return GatewayPayload{
		Op: voice_opcode.Identify,
		D: map[string]interface{}{
			"server_id":  guildId,
			"user_id":    userId,
			"session_id": sessionId,
			"token":      token,
		},
	}
}

func NewVoiceHeartbeat(s int64) GatewayPayload {
	x := strconv.FormatInt(s, 10)
	return GatewayPayload{
		Op: voice_opcode.Heartbeat,
		D:  x,
	}
}

func NewVoiceSelectProtocol(address string, port uint16, mode string) GatewayPayload {
	return GatewayPayload{
		Op: voice_opcode.SelectProtocol,
		D: map[string]interface{}{
			"protocol": "udp",
			"data": map[string]interface{}{
				"address": address,
				"port":    port,
				"mode":    mode,
			},
		},
	}
}

func NewVoiceSpeaking(ssrc uint32) GatewayPayload {
	return GatewayPayload{
		Op: voice_opcode.Speaking,
		D: map[string]interface{}{
			"speaking": 1 << 0,
			"delay":    0,
			"ssrc":     ssrc,
		},
	}
}
