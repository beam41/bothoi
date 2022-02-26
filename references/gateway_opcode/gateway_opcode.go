package gateway_opcode

const (
	// An event was dispatched.
	Dispatch = 0
	// Fired periodically by the client to keep the connection alive.
	Heartbeat = 1
	// Starts a new session during the initial handshake.
	Identify = 2
	// Update the client's presence.
	PresenceUpdate = 3
	// Used to join/leave or move between voice channels.
	VoiceStateUpdate = 4
	// Resume a previous session that was disconnected.
	Resume = 6
	// You should attempt to reconnect and resume immediately.
	Reconnect = 7
	// Request information about offline guild members in a large guild.
	RequestGuildMembers = 8
	// The session has been invalidated. You should reconnect and identify/resume accordingly.
	InvalidSession = 9
	// Sent immediately after connecting, contains the heartbeat_interval to use.
	Hello = 10
	// Sent in response to receiving a heartbeat to acknowledge that it has been received.
	HeartbeatAck = 11
)
