package voice_opcode

const (
	Identify = iota
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
