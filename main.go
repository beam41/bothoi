package main

import (
	"bothoi/command"
	"bothoi/gateway"
	"bothoi/repo"
	"bothoi/voice"
	"log"
)

func main() {
	log.Println("Bot start")
	repo.StartDb()
	voiceClientManager := voice.NewClientManager()
	gatewayClient := gateway.NewClient(voiceClientManager)
	command.NewCommandManager(gatewayClient)
	gatewayClient.Connect()
}
