package models

type SessionState struct {
	V           int8               `json:"v" mapstructure:"v"`
	User        User               `json:"user" mapstructure:"user"`
	Shard       []int16            `json:"shard" mapstructure:"shard"`
	SessionID   string             `json:"session_id" mapstructure:"session_id"`
	Guilds      []UnavailableGuild `json:"guilds" mapstructure:"guilds"`
	Application Application        `json:"application" mapstructure:"application"`
}
