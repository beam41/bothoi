package app_command

import (
	"bothoi/config"
	"bothoi/util/http_util"
	"log"
)

// do register app command
func Register() {
	for _, command := range commandList {
		log.Println("app command request: ", command)
		header := map[string]string{}
		header["Authorization"] = "Bot " + config.BOT_TOKEN

		if (config.DEVELOPMENT) {
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
