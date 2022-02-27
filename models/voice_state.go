package models

type VoiceState struct {
	ChannelID               string      `json:"channel_id" mapstructure:"channel_id"`
	Deaf                    bool        `json:"deaf" mapstructure:"deaf"`
	Mute                    bool        `json:"mute" mapstructure:"mute"`
	RequestToSpeakTimestamp interface{} `json:"request_to_speak_timestamp" mapstructure:"request_to_speak_timestamp"`
	SelfDeaf                bool        `json:"self_deaf" mapstructure:"self_deaf"`
	SelfMute                bool        `json:"self_mute" mapstructure:"self_mute"`
	SelfVideo               bool        `json:"self_video" mapstructure:"self_video"`
	SessionID               string      `json:"session_id" mapstructure:"session_id"`
	Suppress                bool        `json:"suppress" mapstructure:"suppress"`
	UserID                  string      `json:"user_id" mapstructure:"user_id"`
	GuildID                 string      `json:"guild_id" mapstructure:"guild_id"`
}
