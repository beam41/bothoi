package main

import (
	"bothoi/bh_context"
	"bothoi/command"
	"bothoi/gateway"
	"bothoi/voice"
)

func main() {
	gatewayClient := gateway.NewClient()
	voiceClientManager := voice.NewClientManager()
	commandManager := command.NewCommandManager()
	bh_context.SetCtx(gatewayClient, voiceClientManager, commandManager)
	bh_context.GetCommandManager().Register()
	bh_context.GetGatewayClient().Connect()
}
