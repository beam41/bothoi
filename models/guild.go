package models

type Guild struct {
	AfkChannelID                interface{}      `json:"afk_channel_id" mapstructure:"afk_channel_id"`
	AfkTimeout                  int64            `json:"afk_timeout" mapstructure:"afk_timeout"`
	ApplicationCommandCount     int64            `json:"application_command_count" mapstructure:"application_command_count"`
	ApplicationCommandCounts    map[string]int64 `json:"application_command_counts" mapstructure:"application_command_counts"`
	ApplicationID               interface{}      `json:"application_id" mapstructure:"application_id"`
	Banner                      interface{}      `json:"banner" mapstructure:"banner"`
	Channels                    []Channel        `json:"channels" mapstructure:"channels"`
	DefaultMessageNotifications int64            `json:"default_message_notifications" mapstructure:"default_message_notifications"`
	Description                 interface{}      `json:"description" mapstructure:"description"`
	DiscoverySplash             interface{}      `json:"discovery_splash" mapstructure:"discovery_splash"`
	EmbeddedActivities          []interface{}    `json:"embedded_activities" mapstructure:"embedded_activities"`
	Emojis                      []interface{}    `json:"emojis" mapstructure:"emojis"`
	ExplicitContentFilter       int64            `json:"explicit_content_filter" mapstructure:"explicit_content_filter"`
	Features                    []interface{}    `json:"features" mapstructure:"features"`
	GuildHashes                 GuildHashes      `json:"guild_hashes" mapstructure:"guild_hashes"`
	GuildScheduledEvents        []interface{}    `json:"guild_scheduled_events" mapstructure:"guild_scheduled_events"`
	HubType                     interface{}      `json:"hub_type" mapstructure:"hub_type"`
	Icon                        interface{}      `json:"icon" mapstructure:"icon"`
	ID                          string           `json:"id" mapstructure:"id"`
	JoinedAt                    string           `json:"joined_at" mapstructure:"joined_at"`
	Large                       bool             `json:"large" mapstructure:"large"`
	Lazy                        bool             `json:"lazy" mapstructure:"lazy"`
	MaxMembers                  int64            `json:"max_members" mapstructure:"max_members"`
	MaxVideoChannelUsers        int64            `json:"max_video_channel_users" mapstructure:"max_video_channel_users"`
	MemberCount                 int64            `json:"member_count" mapstructure:"member_count"`
	Members                     []GuildMember    `json:"members" mapstructure:"members"`
	MfaLevel                    int64            `json:"mfa_level" mapstructure:"mfa_level"`
	Name                        string           `json:"name" mapstructure:"name"`
	Nsfw                        bool             `json:"nsfw" mapstructure:"nsfw"`
	NsfwLevel                   int64            `json:"nsfw_level" mapstructure:"nsfw_level"`
	OwnerID                     string           `json:"owner_id" mapstructure:"owner_id"`
	PreferredLocale             string           `json:"preferred_locale" mapstructure:"preferred_locale"`
	PremiumProgressBarEnabled   bool             `json:"premium_progress_bar_enabled" mapstructure:"premium_progress_bar_enabled"`
	PremiumSubscriptionCount    int64            `json:"premium_subscription_count" mapstructure:"premium_subscription_count"`
	PremiumTier                 int64            `json:"premium_tier" mapstructure:"premium_tier"`
	Presences                   []interface{}    `json:"presences" mapstructure:"presences"`
	PublicUpdatesChannelID      interface{}      `json:"public_updates_channel_id" mapstructure:"public_updates_channel_id"`
	Region                      string           `json:"region" mapstructure:"region"`
	Roles                       []Role           `json:"roles" mapstructure:"roles"`
	RulesChannelID              interface{}      `json:"rules_channel_id" mapstructure:"rules_channel_id"`
	Splash                      interface{}      `json:"splash" mapstructure:"splash"`
	StageInstances              []interface{}    `json:"stage_instances" mapstructure:"stage_instances"`
	Stickers                    []interface{}    `json:"stickers" mapstructure:"stickers"`
	SystemChannelFlags          int64            `json:"system_channel_flags" mapstructure:"system_channel_flags"`
	SystemChannelID             string           `json:"system_channel_id" mapstructure:"system_channel_id"`
	Threads                     []interface{}    `json:"threads" mapstructure:"threads"`
	Unavailable                 bool             `json:"unavailable" mapstructure:"unavailable"`
	VanityURLCode               interface{}      `json:"vanity_url_code" mapstructure:"vanity_url_code"`
	VerificationLevel           int64            `json:"verification_level" mapstructure:"verification_level"`
	VoiceStates                 []VoiceState     `json:"voice_states" mapstructure:"voice_states"`
}

type GuildHashes struct {
	Channels Hash  `json:"channels" mapstructure:"channels"`
	Metadata Hash  `json:"metadata" mapstructure:"metadata"`
	Roles    Hash  `json:"roles" mapstructure:"roles"`
	Version  int32 `json:"version" mapstructure:"version"`
}

type UnavailableGuild struct {
	Unavailable bool   `json:"unavailable" mapstructure:"unavailable"`
	ID          string `json:"id" mapstructure:"id"`
}
