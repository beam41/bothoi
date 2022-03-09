package commands

import (
	"bothoi/commands/pause"
	"bothoi/commands/play"
	"bothoi/commands/queue"
	"bothoi/config"
	"bothoi/gateway"
	"bothoi/models"
	"bothoi/util/http_util"
	"log"
)

var commandList = []models.AppCommand{
	play.Command,
	queue.Command,
	pause.Command[0],
	pause.Command[1],
}

var executorList = map[string]func(*models.Interaction){
	play.Command.Name:     play.Execute,
	queue.Command.Name:    queue.Execute,
	pause.Command[0].Name: pause.Execute,
	pause.Command[1].Name: pause.Execute,
}

func Register() {
	gateway.SetExecutorList(executorList)
	for _, command := range commandList {
		log.Println("app command request: ", command)
		header := map[string]string{}
		header["Authorization"] = "Bot " + config.BOT_TOKEN

		if config.DEVELOPMENT {
			_, errD := http_util.PostJsonH(config.APP_COMMAND_GUILD_ENDPOINT, command, header)
			if errD != nil {
				log.Println(errD)
				continue
			}
		}
		res, err := http_util.PostJsonH(config.APP_COMMAND_ENDPOINT, command, header)

		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("app command response: ", string(res))
	}
}
