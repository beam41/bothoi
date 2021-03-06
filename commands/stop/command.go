package stop

import (
	"bothoi/models"
	"bothoi/references/app_command_type"
)

var Command = models.AppCommand{
	Type:              app_command_type.ChatInput,
	Name:              "stop",
	Description:       "Stop the player and leave the voice channel",
	DefaultPermission: true,
}
