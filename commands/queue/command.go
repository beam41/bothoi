package queue

import (
	"bothoi/models/discord_models"
	"bothoi/references/app_command_type"
)

var Command = discord_models.AppCommand{
	Type:              app_command_type.ChatInput,
	Name:              "queue",
	Description:       "List the music player queue",
	DefaultPermission: true,
}
