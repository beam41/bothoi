package app_command

import (
	"bothoi/commands/play"
	"bothoi/commands/queue"
	"bothoi/models"

	"github.com/gorilla/websocket"
)

var commandList = []models.AppCommand{
	play.Command,
	queue.Command,
}

var executorList = map[string]func(*models.Interaction, *websocket.Conn){
	play.Command.Name: play.Execute,
	queue.Command.Name: queue.Execute,
}
