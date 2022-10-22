package main

import (
	"bothoi/command"
	"bothoi/command/command_interface"
	"bothoi/gateway"
	"bothoi/gateway/gateway_interface"
	"bothoi/voice"
	"bothoi/voice/voice_interface"
)

func main() {
	var gatewayClient gateway_interface.ClientInterface
	var voiceClientManager voice_interface.ClientManagerInterface
	var commandManager command_interface.CommandManagerInterface
	gatewayClient = gateway.NewClient(commandManager)
	voiceClientManager = voice.NewClientManager(gatewayClient)
	commandManager = command.NewCommandManager(voiceClientManager)
	commandManager.Register()
	gatewayClient.Connect()
}
