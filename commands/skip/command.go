package skip

import (
	"bothoi/models"
	"bothoi/references/app_command_type"
)

var Command = models.AppCommand{
	Type:              app_command_type.ChatInput,
	Name:              "skip",
	Description:       "Skip song in the player",
	DefaultPermission: true,
}
