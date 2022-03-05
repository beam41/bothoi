package app_command

import (
	"bothoi/commands/play"
	"bothoi/commands/queue"
	"bothoi/models"
)

var commandList = []models.AppCommand{
	play.Command,
	queue.Command,
}

var executorList = map[string]func(*models.Interaction){
	play.Command.Name: play.Execute,
	queue.Command.Name: queue.Execute,
}
