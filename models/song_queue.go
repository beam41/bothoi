package models

type SongQueue struct {
	GuildID            string
	CurrVoiceChannelID string
	Songs              []SongItem
}

type SongItem struct {
	YtID        string
	Title       string
	Duration    int64
	RequesterID string
}
