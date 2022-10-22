package bh_context

import (
	"bothoi/command/command_interface"
	"bothoi/gateway/gateway_interface"
	"bothoi/voice/voice_interface"
)

type context struct {
	gatewayClient      gateway_interface.ClientInterface
	voiceClientManager voice_interface.ClientManagerInterface
	commandManager     command_interface.CommandManagerInterface
}

var ctx *context

func SetCtx(gatewayClient gateway_interface.ClientInterface,
	voiceClientManager voice_interface.ClientManagerInterface,
	commandManager command_interface.CommandManagerInterface) {
	ctx = &context{
		gatewayClient,
		voiceClientManager,
		commandManager,
	}
}

func GetGatewayClient() gateway_interface.ClientInterface {
	return ctx.gatewayClient
}

func GetVoiceClientManager() voice_interface.ClientManagerInterface {
	return ctx.voiceClientManager
}

func GetCommandManager() command_interface.CommandManagerInterface {
	return ctx.commandManager
}
