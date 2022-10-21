package discord_models

import "bothoi/models/types"

type VoiceState struct {
	GuildId                 types.Snowflake     `json:"guild_id" mapstructure:"guild_id"`
	ChannelId               *types.Snowflake    `json:"channel_id" mapstructure:"channel_id"`
	UserId                  types.Snowflake     `json:"user_id" mapstructure:"user_id"`
	Member                  GuildMember         `json:"member" mapstructure:"member"`
	SessionId               string              `json:"session_id" mapstructure:"session_id"`
	Deaf                    bool                `json:"deaf" mapstructure:"deaf"`
	Mute                    bool                `json:"mute" mapstructure:"mute"`
	SelfDeaf                bool                `json:"self_deaf" mapstructure:"self_deaf"`
	SelfMute                bool                `json:"self_mute" mapstructure:"self_mute"`
	SelfStream              bool                `json:"self_stream" mapstructure:"self_stream"`
	SelfVideo               bool                `json:"self_video" mapstructure:"self_video"`
	Suppress                bool                `json:"suppress" mapstructure:"suppress"`
	RequestToSpeakTimestamp *types.ISOTimeStamp `json:"request_to_speak_timestamp" mapstructure:"request_to_speak_timestamp"`
}
