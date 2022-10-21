package util

import (
	"bothoi/models/discord_models"
)

func MapInteractionOption(options []discord_models.InteractionOption) map[string]discord_models.InteractionOption {
	mapOp := map[string]discord_models.InteractionOption{}
	for _, option := range options {
		mapOp[option.Name] = option
	}
	return mapOp
}
