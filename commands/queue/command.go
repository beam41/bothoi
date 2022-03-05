package queue

import (
	"bothoi/models"
	"bothoi/references/app_command_type"
)

var Command = models.AppCommand{
	Type:              app_command_type.ChatInput,
	Name:              "queue",
	Description:       "List the music player queue",
	DefaultPermission: true,
}
