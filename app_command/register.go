package app_command

import (
	"bothoi/util/http_util"
	"log"
	"os"
)

// do register app command
func Register() {
	for _, command := range commandList {
		log.Println("app command request: ", command)
		header := make(map[string]string)
		header["Authorization"] = "Bot " + os.Getenv("BOT_TOKEN")

		res, err := http_util.PostJsonH(os.Getenv("APP_COMMAND_ENDPOINT"), command, header)

		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("app command response: ", string(res))
	}
}
