package main

import (
	"bothoi/bh_context"
	"bothoi/command"
	"bothoi/gateway"
	"bothoi/voice"
)

func main() {
	bh_context.Ctx.CommandManager = command.NewCommandManager()
	bh_context.Ctx.GatewayClient = gateway.NewClient()
	bh_context.Ctx.VoiceClientManager = voice.NewClientManager()
	bh_context.Ctx.CommandManager.Register()
	bh_context.Ctx.GatewayClient.Connect()
}
