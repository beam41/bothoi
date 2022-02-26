package app_command

import (
	"bothoi/models"
)

func MapInteraction(data *models.Interaction) {
	executorList[data.Data.Name](data)
}
