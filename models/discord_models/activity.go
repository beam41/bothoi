package discord_models

import "bothoi/models/types"

type Activity struct {
	Name          string              `json:"name" mapstructure:"name"`
	Type          uint8               `json:"type" mapstructure:"type"`
	Url           *string             `json:"url" mapstructure:"url"`
	CreatedAt     types.UnixTimeStamp `json:"created_at" mapstructure:"created_at"`
	Timestamps    ActivityTimestamp   `json:"timestamps" mapstructure:"timestamps"`
	ApplicationId types.Snowflake     `json:"application_id" mapstructure:"application_id"`
	Details       *string             `json:"details" mapstructure:"details"`
	State         *string             `json:"state" mapstructure:"state"`
	Emoji         *ActivityEmoji      `json:"emoji" mapstructure:"emoji"`
	Party         ActivityParty       `json:"party" mapstructure:"party"`
	Assets        ActivityAsset       `json:"assets" mapstructure:"assets"`
	Secrets       ActivitySecret      `json:"secrets" mapstructure:"secrets"`
	Instance      bool                `json:"instance" mapstructure:"instance"`
	Flags         uint16              `json:"flags" mapstructure:"flags"`
	Buttons       []ActivityButton    `json:"buttons" mapstructure:"buttons"`
}

type ActivityTimestamp struct {
	Start types.UnixTimeStamp `json:"start" mapstructure:"start"`
	End   types.UnixTimeStamp `json:"end" mapstructure:"end"`
}

type ActivityEmoji struct {
	Name     string          `json:"name" mapstructure:"name"`
	Id       types.Snowflake `json:"id" mapstructure:"id"`
	Animated bool            `json:"animated" mapstructure:"animated"`
}

type ActivityParty struct {
	Id   string    `json:"id" mapstructure:"id"`
	Size [2]uint32 `json:"size" mapstructure:"size"`
}

type ActivityAsset struct {
	LargeImage string `json:"large_image" mapstructure:"large_image"`
	LargeText  string `json:"large_text" mapstructure:"large_text"`
	SmallImage string `json:"small_image" mapstructure:"small_image"`
	SmallText  string `json:"small_text" mapstructure:"small_text"`
}

type ActivitySecret struct {
	Join     string `json:"join" mapstructure:"join"`
	Spectate string `json:"spectate" mapstructure:"spectate"`
	Match    string `json:"match" mapstructure:"match"`
}

type ActivityButton struct {
	Label string `json:"label" mapstructure:"label"`
	Url   string `json:"url" mapstructure:"url"`
}
