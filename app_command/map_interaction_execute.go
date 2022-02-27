package app_command

import (
	"bothoi/models"
)

func MapInteractionExecute(data *models.Interaction) {
	executorList[data.Data.Name](data)
}
