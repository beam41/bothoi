package models

type SessionState struct {
	V                    int8               `json:"v" mapstructure:"v"`
	UserSettings         interface{}        `json:"user_settings" mapstructure:"user_settings"`
	User                 User               `json:"user" mapstructure:"user"`
	Shard                []int16            `json:"shard" mapstructure:"shard"`
	SessionID            string             `json:"session_id" mapstructure:"session_id"`
	Relationships        []interface{}      `json:"relationships" mapstructure:"relationships"`
	PrivateChannels      []interface{}      `json:"private_channels" mapstructure:"private_channels"`
	Presences            []interface{}      `json:"presences" mapstructure:"presences"`
	Guilds               []UnavailableGuild `json:"guilds" mapstructure:"guilds"`
	GuildJoinRequests    []interface{}      `json:"guild_join_requests" mapstructure:"guild_join_requests"`
	GeoOrderedRTCRegions []string           `json:"geo_ordered_rtc_regions" mapstructure:"geo_ordered_rtc_regions"`
	Application          Application        `json:"application" mapstructure:"application"`
	Trace                []string           `json:"_trace" mapstructure:"_trace"`
}
