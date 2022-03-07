package models

type SongQueue struct {
	GuildID        string
	VoiceChannelID string
	SessionID      *string
	Songs          []SongItem
	VoiceServer    *VoiceServer

}

type SongItem struct {
	YtID        string
	Title       string
	Duration    int64
	RequesterID string
}

type VoiceServer struct {
	Token    string `json:"token" mapstructure:"token"`
	GuildID  string `json:"guild_id" mapstructure:"guild_id"`
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`
}
