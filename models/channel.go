package models

type Channel struct {
	ID                   string  `json:"id" mapstructure:"id"`
	Name                 string  `json:"name" mapstructure:"name"`
	PermissionOverwrites []any   `json:"permission_overwrites" mapstructure:"permission_overwrites"`
	Position             int64   `json:"position" mapstructure:"position"`
	Type                 int64   `json:"type" mapstructure:"type"`
	LastMessageID        *string `json:"last_message_id" mapstructure:"last_message_id"`
	ParentID             *string `json:"parent_id,omitempty" mapstructure:"parent_id"`
	RateLimitPerUser     *int64  `json:"rate_limit_per_user,omitempty" mapstructure:"rate_limit_per_user"`
	Topic                any     `json:"topic" mapstructure:"topic"`
	Bitrate              *int64  `json:"bitrate,omitempty" mapstructure:"bitrate"`
	RTCRegion            any     `json:"rtc_region" mapstructure:"rtc_region"`
	UserLimit            *int64  `json:"user_limit,omitempty" mapstructure:"user_limit"`
	Nsfw                 *bool   `json:"nsfw,omitempty" mapstructure:"nsfw"`
}
