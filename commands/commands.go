package commands

import (
	"bothoi/commands/pause"
	"bothoi/commands/play"
	"bothoi/commands/queue"
	"bothoi/commands/skip"
	"bothoi/commands/stop"
	"bothoi/config"
	"bothoi/models"
	"bothoi/util/http_util"
	"log"
)

var ExecutorList = map[string]func(*models.Interaction){
	play.Command.Name:     play.Execute,
	queue.Command.Name:    queue.Execute,
	pause.Command[0].Name: pause.Execute,
	pause.Command[1].Name: pause.Execute,
	stop.Command.Name:     stop.Execute,
	skip.Command.Name:     skip.Execute,
}

func Register() {
	if config.NO_COMMAND_REGISTER {
		return
	}
	var commandList = []models.AppCommand{
		play.Command,
		queue.Command,
		pause.Command[0],
		pause.Command[1],
		stop.Command,
		skip.Command,
	}
	for _, command := range commandList {
		log.Println("app command request: ", command)
		header := map[string]string{}
		header["Authorization"] = "Bot " + config.BOT_TOKEN

		if config.DEVELOPMENT {
			res, errD := http_util.PostJsonH(config.APP_COMMAND_GUILD_ENDPOINT, command, header)
			log.Println("guild command response: ", string(res))
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
