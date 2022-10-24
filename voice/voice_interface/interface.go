package voice_interface

import (
	"bothoi/models/types"
)

type ClientManagerInterface interface {
	PauseClient(guildId types.Snowflake) (bool, error)
	SkipSong(guildId types.Snowflake) error
	StartClient(guildId, channelId types.Snowflake) error
	StopClient(guildId types.Snowflake) error
}
