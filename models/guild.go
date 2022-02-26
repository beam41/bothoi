package models

type Guild struct {
	AfkChannelID                interface{}      `json:"afk_channel_id"`
	AfkTimeout                  int64            `json:"afk_timeout"`
	ApplicationCommandCount     int64            `json:"application_command_count"`
	ApplicationCommandCounts    map[string]int64 `json:"application_command_counts"`
	ApplicationID               interface{}      `json:"application_id"`
	Banner                      interface{}      `json:"banner"`
	Channels                    []Channel        `json:"channels"`
	DefaultMessageNotifications int64            `json:"default_message_notifications"`
	Description                 interface{}      `json:"description"`
	DiscoverySplash             interface{}      `json:"discovery_splash"`
	EmbeddedActivities          []interface{}    `json:"embedded_activities"`
	Emojis                      []interface{}    `json:"emojis"`
	ExplicitContentFilter       int64            `json:"explicit_content_filter"`
	Features                    []interface{}    `json:"features"`
	GuildHashes                 GuildHashes      `json:"guild_hashes"`
	GuildScheduledEvents        []interface{}    `json:"guild_scheduled_events"`
	HubType                     interface{}      `json:"hub_type"`
	Icon                        interface{}      `json:"icon"`
	ID                          string           `json:"id"`
	JoinedAt                    string           `json:"joined_at"`
	Large                       bool             `json:"large"`
	Lazy                        bool             `json:"lazy"`
	MaxMembers                  int64            `json:"max_members"`
	MaxVideoChannelUsers        int64            `json:"max_video_channel_users"`
	MemberCount                 int64            `json:"member_count"`
	Members                     []Member         `json:"members"`
	MfaLevel                    int64            `json:"mfa_level"`
	Name                        string           `json:"name"`
	Nsfw                        bool             `json:"nsfw"`
	NsfwLevel                   int64            `json:"nsfw_level"`
	OwnerID                     string           `json:"owner_id"`
	PreferredLocale             string           `json:"preferred_locale"`
	PremiumProgressBarEnabled   bool             `json:"premium_progress_bar_enabled"`
	PremiumSubscriptionCount    int64            `json:"premium_subscription_count"`
	PremiumTier                 int64            `json:"premium_tier"`
	Presences                   []interface{}    `json:"presences"`
	PublicUpdatesChannelID      interface{}      `json:"public_updates_channel_id"`
	Region                      string           `json:"region"`
	Roles                       []Role           `json:"roles"`
	RulesChannelID              interface{}      `json:"rules_channel_id"`
	Splash                      interface{}      `json:"splash"`
	StageInstances              []interface{}    `json:"stage_instances"`
	Stickers                    []interface{}    `json:"stickers"`
	SystemChannelFlags          int64            `json:"system_channel_flags"`
	SystemChannelID             string           `json:"system_channel_id"`
	Threads                     []interface{}    `json:"threads"`
	Unavailable                 bool             `json:"unavailable"`
	VanityURLCode               interface{}      `json:"vanity_url_code"`
	VerificationLevel           int64            `json:"verification_level"`
	VoiceStates                 []VoiceState     `json:"voice_states"`
}

type Channel struct {
	ID                   string        `json:"id"`
	Name                 string        `json:"name"`
	PermissionOverwrites []interface{} `json:"permission_overwrites"`
	Position             int64         `json:"position"`
	Type                 int64         `json:"type"`
	LastMessageID        *string       `json:"last_message_id"`
	ParentID             *string       `json:"parent_id,omitempty"`
	RateLimitPerUser     *int64        `json:"rate_limit_per_user,omitempty"`
	Topic                interface{}   `json:"topic"`
	Bitrate              *int64        `json:"bitrate,omitempty"`
	RTCRegion            interface{}   `json:"rtc_region"`
	UserLimit            *int64        `json:"user_limit,omitempty"`
	Nsfw                 *bool         `json:"nsfw,omitempty"`
}

type GuildHashes struct {
	Channels Channels `json:"channels"`
	Metadata Channels `json:"metadata"`
	Roles    Channels `json:"roles"`
	Version  int64    `json:"version"`
}

type Channels struct {
	Hash    string `json:"hash"`
	Omitted bool   `json:"omitted"`
}

type Member struct {
	Deaf                       bool        `json:"deaf"`
	HoistedRole                interface{} `json:"hoisted_role"`
	JoinedAt                   string      `json:"joined_at"`
	Mute                       bool        `json:"mute"`
	Roles                      []string    `json:"roles"`
	User                       User        `json:"user"`
	Avatar                     interface{} `json:"avatar"`
	CommunicationDisabledUntil interface{} `json:"communication_disabled_until"`
	Nick                       interface{} `json:"nick"`
	Pending                    *bool       `json:"pending,omitempty"`
	PremiumSince               interface{} `json:"premium_since"`
}

type User struct {
	Avatar        *string `json:"avatar"`
	Discriminator string  `json:"discriminator"`
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Bot           *bool   `json:"bot,omitempty"`
}

type Role struct {
	Color        int64       `json:"color"`
	Hoist        bool        `json:"hoist"`
	Icon         interface{} `json:"icon"`
	ID           string      `json:"id"`
	Managed      bool        `json:"managed"`
	Mentionable  bool        `json:"mentionable"`
	Name         string      `json:"name"`
	Permissions  string      `json:"permissions"`
	Position     int64       `json:"position"`
	UnicodeEmoji interface{} `json:"unicode_emoji"`
	Tags         *Tags       `json:"tags,omitempty"`
}

type Tags struct {
	BotID string `json:"bot_id"`
}

type VoiceState struct {
	ChannelID               string      `json:"channel_id"`
	Deaf                    bool        `json:"deaf"`
	Mute                    bool        `json:"mute"`
	RequestToSpeakTimestamp interface{} `json:"request_to_speak_timestamp"`
	SelfDeaf                bool        `json:"self_deaf"`
	SelfMute                bool        `json:"self_mute"`
	SelfVideo               bool        `json:"self_video"`
	SessionID               string      `json:"session_id"`
	Suppress                bool        `json:"suppress"`
	UserID                  string      `json:"user_id"`
}
