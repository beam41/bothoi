package models

type Role struct {
	Color        int64     `json:"color" mapstructure:"color"`
	Hoist        bool      `json:"hoist" mapstructure:"hoist"`
	Icon         any       `json:"icon" mapstructure:"icon"`
	ID           string    `json:"id" mapstructure:"id"`
	Managed      bool      `json:"managed" mapstructure:"managed"`
	Mentionable  bool      `json:"mentionable" mapstructure:"mentionable"`
	Name         string    `json:"name" mapstructure:"name"`
	Permissions  string    `json:"permissions" mapstructure:"permissions"`
	Position     int64     `json:"position" mapstructure:"position"`
	UnicodeEmoji any       `json:"unicode_emoji" mapstructure:"unicode_emoji"`
	Tags         *RoleTags `json:"tags,omitempty" mapstructure:"tags"`
}

type RoleTags struct {
	BotID string `json:"bot_id" mapstructure:"bot_id"`
}
