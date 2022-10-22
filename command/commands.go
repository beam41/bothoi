package command

import (
	"bothoi/command/command_interface"
	"bothoi/models/discord_models"
)

type commandManager struct {
	executorList map[string]func(*discord_models.Interaction)
}

func NewCommandManager() command_interface.CommandManagerInterface {
	return &commandManager{
		executorList: map[string]func(*discord_models.Interaction){
			commandPlay:   executePlay,
			commandQueue:  executeQueue,
			commandPause0: executePause,
			commandPause1: executePause,
			commandStop:   executeStop,
			commandSkip:   executeSkip,
		},
	}
}

func (cm *commandManager) MapInteractionExecute(data *discord_models.Interaction) {
	if interaction, ok := cm.executorList[data.Data.Name]; ok {
		interaction(data)
	}
}
