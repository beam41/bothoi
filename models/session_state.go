package models

type SessionState struct {
	V                    int8                     `json:"v"`
	UserSettings         SessionStateUserSettings `json:"user_settings"`
	User                 SessionStateUser         `json:"user"`
	Shard                []int16                  `json:"shard"`
	SessionID            string                   `json:"session_id"`
	Relationships        []interface{}            `json:"relationships"`
	PrivateChannels      []interface{}            `json:"private_channels"`
	Presences            []interface{}            `json:"presences"`
	Guilds               []SessionStateGuild      `json:"guilds"`
	GuildJoinRequests    []interface{}            `json:"guild_join_requests"`
	GeoOrderedRTCRegions []string                 `json:"geo_ordered_rtc_regions"`
	Application          SessionStateApplication  `json:"application"`
	Trace                []string                 `json:"_trace"`
}

type SessionStateApplication struct {
	ID    string `json:"id"`
	Flags int64  `json:"flags"`
}

type SessionStateGuild struct {
	Unavailable bool   `json:"unavailable"`
	ID          string `json:"id"`
}

type SessionStateUser struct {
	Verified      bool   `json:"verified"`
	Username      string `json:"username"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	ID            string `json:"id"`
	Flags         int32  `json:"flags"`
	Email         string `json:"email"`
	Discriminator string `json:"discriminator"`
	Bot           bool   `json:"bot"`
	Avatar        string `json:"avatar"`
}

type SessionStateUserSettings struct {
}
