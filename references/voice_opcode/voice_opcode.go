package voice_opcode

type VoiceOpcode byte

const (
	Identify VoiceOpcode = iota
	SelectProtocol
	Ready
	Heartbeat
	SessionDescription
	Speaking
	HeartbeatAck
	Resume
	Hello
	Resumed
	_
	_
	_
	ClientDisconnect
)
