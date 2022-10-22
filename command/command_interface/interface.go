package command_interface

import "bothoi/models/discord_models"

type CommandManagerInterface interface {
	MapInteractionExecute(data *discord_models.Interaction)
	Register()
}
