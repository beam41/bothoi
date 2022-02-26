package models

type Interaction struct {
	ApplicationID string            `json:"application_id"`
	ChannelID     string            `json:"channel_id"`
	Data          InteractionData   `json:"data"`
	GuildID       string            `json:"guild_id"`
	GuildLocale   string            `json:"guild_locale"`
	ID            string            `json:"id"`
	Locale        string            `json:"locale"`
	Member        InteractionMember `json:"member"`
	Token         string            `json:"token"`
	Type          int64             `json:"type"`
	Version       int64             `json:"version"`
}

type InteractionData struct {
	ID      string              `json:"id"`
	Name    string              `json:"name"`
	Options []InteractionOption `json:"options"`
	Type    int64               `json:"type"`
}

type InteractionOption struct {
	Name  string      `json:"name"`
	Type  int64       `json:"type"`
	Value interface{} `json:"value"`
}

type InteractionMember struct {
	Avatar                     interface{}     `json:"avatar"`
	CommunicationDisabledUntil interface{}     `json:"communication_disabled_until"`
	Deaf                       bool            `json:"deaf"`
	IsPending                  bool            `json:"is_pending"`
	JoinedAt                   string          `json:"joined_at"`
	Mute                       bool            `json:"mute"`
	Nick                       interface{}     `json:"nick"`
	Pending                    bool            `json:"pending"`
	Permissions                string          `json:"permissions"`
	PremiumSince               interface{}     `json:"premium_since"`
	Roles                      []interface{}   `json:"roles"`
	User                       InteractionUser `json:"user"`
}

type InteractionUser struct {
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
	PublicFlags   int64  `json:"public_flags"`
	Username      string `json:"username"`
}
