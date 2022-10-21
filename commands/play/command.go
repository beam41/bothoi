package play

import (
	"bothoi/models/discord_models"
	"bothoi/references/app_command_option_type"
	"bothoi/references/app_command_type"
)

var Command = discord_models.AppCommand{
	Type:              app_command_type.ChatInput,
	Name:              "play",
	Description:       "Play a song",
	DefaultPermission: true,
	Options: []discord_models.AppCommandOption{
		{
			Type:        app_command_option_type.String,
			Name:        "song",
			Description: "The song to play",
			Required:    true,
		},
	},
}
