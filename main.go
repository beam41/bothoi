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
	gatewayClient := gateway.NewClient()
	voiceClientManager := voice.NewClientManager(gatewayClient)
	command.NewCommandManager(gatewayClient, voiceClientManager)
	gatewayClient.Connect()
}
