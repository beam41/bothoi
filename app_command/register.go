package app_command

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// do register app command
func Register() {
	for _, command := range commandList {

		payload := new(bytes.Buffer)
		json.NewEncoder(payload).Encode(command)
		req, err := http.NewRequest("POST", os.Getenv("APP_COMMAND_ENDPOINT"), payload)

		if err != nil {
			log.Println(err)
			continue
		}
		req.Header.Add("Authorization", "Bot "+os.Getenv("BOT_TOKEN"))
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("app command response: ", string(body))
		}
	}
}
