package command

import (
	"bothoi/command/command_interface"
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/util/http_util"
	"log"
)

type commandManager struct {
	executorList map[string]func(*discord_models.Interaction)
}

func NewCommandManager() command_interface.CommandManagerInterface {
	return &commandManager{
		executorList: map[string]func(*discord_models.Interaction){
			commandPlay.Name:     executePlay,
			commandQueue.Name:    executeQueue,
			commandPause[0].Name: executePause,
			commandPause[1].Name: executePause,
			commandStop.Name:     executeStop,
			commandSkip.Name:     executeSkip,
		},
	}
}

func (cm *commandManager) MapInteractionExecute(data *discord_models.Interaction) {
	if interaction, ok := cm.executorList[data.Data.Name]; ok {
		interaction(data)
	}
}

func (cm *commandManager) Register() {
	if config.NoCommandRegister {
		return
	}

	var commandList = []discord_models.AppCommand{
		commandPlay,
		commandQueue,
		commandPause[0],
		commandPause[1],
		commandStop,
		commandSkip,
	}

	for _, command := range commandList {
		log.Println("app command request: ", command)
		header := map[string]string{}
		header["Authorization"] = "Bot " + config.BotToken

		if config.Development {
			res, errD := http_util.PostJsonH(config.AppCommandEndpoint, command, header)
			log.Println("guild command response: ", string(res))
			if errD != nil {
				log.Println(errD)
				continue
			}
		}
		res, err := http_util.PostJsonH(config.AppCommandEndpoint, command, header)

		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("app command response: ", string(res))
	}
}
