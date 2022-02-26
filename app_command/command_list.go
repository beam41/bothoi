package app_command

import (
	"bothoi/commands/play"
	"bothoi/models"
)

var commandList = []models.AppCommand{
	play.Command,
}

var executorList = map[string]func(*models.Interaction){
	play.Command.Name: play.Execute,
}
