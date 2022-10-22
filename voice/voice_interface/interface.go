package voice_interface

import (
	"bothoi/models"
	"bothoi/models/types"
)

type ClientManagerInterface interface {
	AppendSongToSongQueue(guildId types.Snowflake, songItem models.SongItem) int
	GetSongQueue(guildId types.Snowflake, start, end int) (playing bool, queue []models.SongItem)
	PauseClient(guildId types.Snowflake) (bool, error)
	SkipSong(guildId types.Snowflake) error
	GetVoiceChannelId(guildId types.Snowflake) types.Snowflake
	StartClient(guildId, channelId types.Snowflake) error
	StopClient(guildId types.Snowflake) error
}
