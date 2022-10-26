package voice_interface

import (
	"bothoi/models/types"
)

type ClientManagerInterface interface {
	PauseClient(guildID types.Snowflake) (bool, error)
	SkipSong(guildID types.Snowflake) error
	StartClient(guildID, channelID types.Snowflake) error
	StopClient(guildID types.Snowflake) error
}
