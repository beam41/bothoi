package main

import (
	"bothoi/bh_context"
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
	voiceClientManager := voice.NewClientManager()
	commandManager := command.NewCommandManager()
	bh_context.SetCtx(gatewayClient, voiceClientManager, commandManager)
	gatewayClient.Connect()
}
