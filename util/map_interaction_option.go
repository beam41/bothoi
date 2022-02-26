package util

import "bothoi/models"

func MapInteractionOption(options []models.InteractionOption) map[string]models.InteractionOption {
	mapOp := make(map[string]models.InteractionOption)
	for _, option := range options {
		mapOp[option.Name] = option
	}
	return mapOp
}
