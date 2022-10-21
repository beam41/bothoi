package stop

import (
	"bothoi/models/discord_models"
	"bothoi/references/app_command_type"
)

var Command = discord_models.AppCommand{
	Type:              app_command_type.ChatInput,
	Name:              "stop",
	Description:       "Stop the player and leave the voice channel",
	DefaultPermission: true,
}
