package discord_models

import (
	"bothoi/config"
	"bothoi/models/types"
	"bothoi/references/gateway_opcode"
	"bothoi/references/voice_opcode"
	"strconv"
)

type GatewayPayload struct {
	Op gateway_opcode.GatewayOpcode `json:"op" mapstructure:"op"`
	S  *uint64                      `json:"s,omitempty" mapstructure:"s"`
	T  string                       `json:"t" mapstructure:"t"`
	D  any                          `json:"d" mapstructure:"d"`
}

type VoiceGatewayPayload struct {
	Op voice_opcode.VoiceOpcode `json:"op" mapstructure:"op"`
	S  *uint64                  `json:"s,omitempty" mapstructure:"s"`
	T  string                   `json:"t" mapstructure:"t"`
	D  any                      `json:"d" mapstructure:"d"`
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
		D: map[string]any{
			"token":   config.BotToken,
			"intents": config.GatewayIntent,
			"properties": map[string]any{
				"$os":      "linux",
				"$browser": "bothoi",
				"$device":  "bothoi",
			},
			"compress": false,
		},
	}
}

func NewResume(s *uint64, sessionID string) GatewayPayload {
	return GatewayPayload{
		Op: gateway_opcode.Resume,
		D: map[string]any{
			"token":      config.BotToken,
			"session_id": sessionID,
			"seq":        s,
		},
	}
}

func NewVoiceStateUpdate(guildID types.Snowflake, voiceID *types.Snowflake, mute, deaf bool) GatewayPayload {
	return GatewayPayload{
		Op: gateway_opcode.VoiceStateUpdate,
		D: map[string]any{
			"guild_id":   guildID,
			"channel_id": voiceID,
			"self_mute":  mute,
			"self_deaf":  deaf,
		},
	}
}

func NewVoiceIdentify(guildID, userID types.Snowflake, sessionID, token string) VoiceGatewayPayload {
	return VoiceGatewayPayload{
		Op: voice_opcode.Identify,
		D: map[string]any{
			"server_id":  guildID,
			"user_id":    userID,
			"session_id": sessionID,
			"token":      token,
		},
	}
}

func NewVoiceResume(guildID types.Snowflake, sessionID, token string) VoiceGatewayPayload {
	return VoiceGatewayPayload{
		Op: voice_opcode.Resume,
		D: map[string]any{
			"server_id":  guildID,
			"session_id": sessionID,
			"token":      token,
		},
	}
}

func NewVoiceHeartbeat(s int64) VoiceGatewayPayload {
	x := strconv.FormatInt(s, 10)
	return VoiceGatewayPayload{
		Op: voice_opcode.Heartbeat,
		D:  x,
	}
}

func NewVoiceSelectProtocol(address string, port uint16, mode string) VoiceGatewayPayload {
	return VoiceGatewayPayload{
		Op: voice_opcode.SelectProtocol,
		D: map[string]any{
			"protocol": "udp",
			"data": map[string]any{
				"address": address,
				"port":    port,
				"mode":    mode,
			},
		},
	}
}

func NewVoiceSpeaking(ssrc uint32) VoiceGatewayPayload {
	return VoiceGatewayPayload{
		Op: voice_opcode.Speaking,
		D: map[string]any{
			"speaking": 1 << 0,
			"delay":    0,
			"ssrc":     ssrc,
		},
	}
}
