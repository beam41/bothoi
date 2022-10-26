package discord_models

import "bothoi/models/types"

type Channel struct {
	ID                            types.Snowflake     `json:"id,string" mapstructure:"id"`
	Type                          uint8               `json:"type" mapstructure:"type"`
	GuildID                       types.Snowflake     `json:"guild_id,string" mapstructure:"guild_id"`
	Position                      uint16              `json:"position" mapstructure:"position"`
	PermissionOverwrites          []OverWrite         `json:"permission_overwrites" mapstructure:"permission_overwrites"`
	Name                          *string             `json:"name" mapstructure:"name"`
	Topic                         *string             `json:"topic" mapstructure:"topic"`
	Nsfw                          bool                `json:"nsfw" mapstructure:"nsfw"`
	LastMessageID                 *types.Snowflake    `json:"last_message_id,string" mapstructure:"last_message_id"`
	Bitrate                       uint32              `json:"bitrate" mapstructure:"bitrate"`
	UserLimit                     uint32              `json:"user_limit" mapstructure:"user_limit"`
	RateLimitPerUser              uint16              `json:"rate_limit_per_user" mapstructure:"rate_limit_per_user"`
	Recipients                    []User              `json:"recipients" mapstructure:"recipients"`
	Icon                          *string             `json:"icon" mapstructure:"icon"`
	OwnerID                       types.Snowflake     `json:"owner_id,string" mapstructure:"owner_id"`
	ApplicationID                 types.Snowflake     `json:"application_id,string" mapstructure:"application_id"`
	ParentID                      *types.Snowflake    `json:"parent_id,string" mapstructure:"parent_id"`
	LastPinTimestamp              *types.ISOTimeStamp `json:"last_pin_timestamp" mapstructure:"last_pin_timestamp"`
	RtcRegion                     *string             `json:"rtc_region" mapstructure:"rtc_region"`
	VideoQualityMode              uint8               `json:"video_quality_mode" mapstructure:"video_quality_mode"`
	MessageCount                  uint64              `json:"message_count" mapstructure:"message_count"`
	MemberCount                   uint32              `json:"member_count" mapstructure:"member_count"`
	ThreadMetadata                ThreadMetadata      `json:"thread_metadata" mapstructure:"thread_metadata"`
	Member                        *ThreadMember       `json:"member" mapstructure:"member"`
	DefaultAutoArchiveDuration    uint16              `json:"default_auto_archive_duration" mapstructure:"default_auto_archive_duration"`
	Permissions                   string              `json:"permissions" mapstructure:"permissions"`
	Flags                         uint8               `json:"flags" mapstructure:"flags"`
	TotalMessageSent              uint64              `json:"total_message_sent" mapstructure:"total_message_sent"`
	AvailableTags                 []Tag               `json:"available_tags" mapstructure:"available_tags"`
	AppliedTags                   []types.Snowflake   `json:"applied_tags,string" mapstructure:"applied_tags"`
	DefaultReactionEmoji          *DefaultReaction    `json:"default_reaction_emoji" mapstructure:"default_reaction_emoji"`
	DefaultThreadRateLimitPerUser uint16              `json:"default_thread_rate_limit_per_user" mapstructure:"default_thread_rate_limit_per_user"`
	DefaultSortOrder              *uint16             `json:"default_sort_order" mapstructure:"default_sort_order"`
}

type OverWrite struct {
	ID    types.Snowflake `json:"id,string" mapstructure:"id"`
	Type  byte            `json:"type" mapstructure:"type"`
	Allow string          `json:"allow" mapstructure:"allow"`
	Deny  string          `json:"deny" mapstructure:"deny"`
}

type ThreadMetadata struct {
	Archived            bool                `json:"archived" mapstructure:"archived"`
	AutoArchiveDuration uint16              `json:"auto_archive_duration" mapstructure:"auto_archive_duration"`
	ArchiveTimestamp    types.ISOTimeStamp  `json:"archive_timestamp" mapstructure:"archive_timestamp"`
	Locked              bool                `json:"locked" mapstructure:"locked"`
	Invitable           bool                `json:"invitable" mapstructure:"invitable"`
	CreateTimestamp     *types.ISOTimeStamp `json:"create_timestamp" mapstructure:"create_timestamp"`
}

type ThreadMember struct {
	ID            types.Snowflake    `json:"id,string" mapstructure:"id"`
	UserID        types.Snowflake    `json:"user_id,string" mapstructure:"user_id"`
	JoinTimestamp types.ISOTimeStamp `json:"join_timestamp" mapstructure:"join_timestamp"`
	Flags         uint64             `json:"flags" mapstructure:"flags"`
}

type Tag struct {
	ID        types.Snowflake  `json:"id,string" mapstructure:"id"`
	Name      string           `json:"name" mapstructure:"name"`
	Moderated bool             `json:"moderated" mapstructure:"moderated"`
	EmojiID   *types.Snowflake `json:"emoji_id,string" mapstructure:"emoji_id"`
	EmojiName *string          `json:"emoji_name" mapstructure:"emoji_name"`
}

type DefaultReaction struct {
	EmojiID   *types.Snowflake `json:"emoji_id,string" mapstructure:"emoji_id"`
	EmojiName *string          `json:"emoji_name" mapstructure:"emoji_name"`
}
