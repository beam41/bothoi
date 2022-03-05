package app_command

import (
	"bothoi/models"

	"github.com/gorilla/websocket"
)

func MapInteractionExecute(data *models.Interaction, c *websocket.Conn) {
	executorList[data.Data.Name](data, c)
}
