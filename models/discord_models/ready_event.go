package discord_models

type ReadyEvent struct {
	V           uint8              `json:"v" mapstructure:"v"`
	User        User               `json:"user" mapstructure:"user"`
	Shard       [2]uint16          `json:"shard" mapstructure:"shard"`
	SessionId   string             `json:"session_id" mapstructure:"session_id"`
	Guilds      []UnavailableGuild `json:"guilds" mapstructure:"guilds"`
	Application Application        `json:"application" mapstructure:"application"`
}
