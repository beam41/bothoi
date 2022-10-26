package discord_models

import (
	"bothoi/models/types"
)

type Guild struct {
	ID                          types.Snowflake  `json:"id,string" mapstructure:"id"`
	Name                        string           `json:"name" mapstructure:"name"`
	Icon                        *string          `json:"icon" mapstructure:"icon"`
	IconHash                    *string          `json:"icon_hash" mapstructure:"icon_hash"`
	Splash                      *string          `json:"splash" mapstructure:"splash"`
	DiscoverySplash             *string          `json:"discovery_splash" mapstructure:"discovery_splash"`
	Owner                       bool             `json:"owner" mapstructure:"owner"`
	OwnerID                     types.Snowflake  `json:"owner_id,string" mapstructure:"owner_id"`
	Permissions                 string           `json:"permissions" mapstructure:"permissions"`
	AfkChannelID                *types.Snowflake `json:"afk_channel_id,string" mapstructure:"afk_channel_id"`
	AfkTimeout                  uint16           `json:"afk_timeout" mapstructure:"afk_timeout"`
	WidgetEnabled               bool             `json:"widget_enabled" mapstructure:"widget_enabled"`
	WidgetChannelID             *types.Snowflake `json:"widget_channel_id,string" mapstructure:"widget_channel_id"`
	VerificationLevel           uint8            `json:"verification_level" mapstructure:"verification_level"`
	DefaultMessageNotifications uint8            `json:"default_message_notifications" mapstructure:"default_message_notifications"`
	ExplicitContentFilter       uint8            `json:"explicit_content_filter" mapstructure:"explicit_content_filter"`
	Roles                       []Role           `json:"roles" mapstructure:"roles"`
	Emojis                      []Emoji          `json:"emojis" mapstructure:"emojis"`
	Features                    []string         `json:"features" mapstructure:"features"`
	MfaLevel                    uint8            `json:"mfa_level" mapstructure:"mfa_level"`
	ApplicationID               *types.Snowflake `json:"application_id,string" mapstructure:"application_id"`
	SystemChannelID             *types.Snowflake `json:"system_channel_id,string" mapstructure:"system_channel_id"`
	SystemChannelFlags          uint8            `json:"system_channel_flags" mapstructure:"system_channel_flags"`
	RulesChannelID              *types.Snowflake `json:"rules_channel_id,string" mapstructure:"rules_channel_id"`
	MaxPresences                *uint32          `json:"max_presences" mapstructure:"max_presences"`
	MaxMembers                  uint32           `json:"max_members" mapstructure:"max_members"`
	VanityURLCode               *string          `json:"vanity_url_code" mapstructure:"vanity_url_code"`
	Description                 *string          `json:"description" mapstructure:"description"`
	Banner                      *string          `json:"banner" mapstructure:"banner"`
	PremiumTier                 uint8            `json:"premium_tier" mapstructure:"premium_tier"`
	PremiumSubscriptionCount    uint32           `json:"premium_subscription_count" mapstructure:"premium_subscription_count"`
	PreferredLocale             string           `json:"preferred_locale" mapstructure:"preferred_locale"`
	PublicUpdatesChannelID      *types.Snowflake `json:"public_updates_channel_id,string" mapstructure:"public_updates_channel_id"`
	MaxVideoChannelUsers        uint32           `json:"max_video_channel_users" mapstructure:"max_video_channel_users"`
	ApproximateMemberCount      uint32           `json:"approximate_member_count" mapstructure:"approximate_member_count"`
	ApproximatePresenceCount    uint32           `json:"approximate_presence_count" mapstructure:"approximate_presence_count"`
	WelcomeScreen               WelcomeScreen    `json:"welcome_screen" mapstructure:"welcome_screen"`
	NsfwLevel                   uint8            `json:"nsfw_level" mapstructure:"nsfw_level"`
	Stickers                    []Sticker        `json:"stickers" mapstructure:"stickers"`
	PremiumProgressBarEnabled   bool             `json:"premium_progress_bar_enabled" mapstructure:"premium_progress_bar_enabled"`
}

type GuildCreate struct {
	Guild                `mapstructure:",squash"`
	JoinedAt             types.ISOTimeStamp    `json:"joined_at" mapstructure:"joined_at"`
	Large                bool                  `json:"large" mapstructure:"large"`
	Unavailable          bool                  `json:"unavailable" mapstructure:"unavailable"`
	MemberCount          uint32                `json:"member_count" mapstructure:"member_count"`
	VoiceStates          []VoiceState          `json:"voice_states" mapstructure:"voice_states"`
	Members              []GuildMember         `json:"members" mapstructure:"members"`
	Channels             []Channel             `json:"channels" mapstructure:"channels"`
	Threads              []Channel             `json:"threads" mapstructure:"threads"`
	Presences            []PresenceUpdate      `json:"presences" mapstructure:"presences"`
	StageInstances       []StageInstance       `json:"stage_instances" mapstructure:"stage_instances"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events" mapstructure:"guild_scheduled_events"`
}

type UnavailableGuild struct {
	Unavailable bool            `json:"unavailable" mapstructure:"unavailable"`
	ID          types.Snowflake `json:"id,string" mapstructure:"id"`
}
