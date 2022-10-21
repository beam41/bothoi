package pause

import (
	"bothoi/models/discord_models"
	"bothoi/references/app_command_type"
)

var Command = []discord_models.AppCommand{
	{
		Type:              app_command_type.ChatInput,
		Name:              "pause",
		Description:       "Pause/Resume the player",
		DefaultPermission: true,
	},
	{
		Type:              app_command_type.ChatInput,
		Name:              "resume",
		Description:       "Pause/Resume the player",
		DefaultPermission: true,
	},
}
