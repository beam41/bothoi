package main

import (
	"bothoi/app_command"
	"bothoi/gateway"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	app_command.Register()
	gateway.Connect()
}
