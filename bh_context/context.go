package bh_context

import (
	"bothoi/command/command_interface"
	"bothoi/gateway/gateway_interface"
	"bothoi/voice/voice_interface"
)

type Context struct {
	GatewayClient      gateway_interface.ClientInterface
	CommandManager     command_interface.CommandManagerInterface
	VoiceClientManager voice_interface.ClientManagerInterface
}

var Ctx = &Context{
	GatewayClient:      nil,
	CommandManager:     nil,
	VoiceClientManager: nil,
}
