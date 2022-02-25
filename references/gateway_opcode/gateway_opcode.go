package gateway_opcode

const (
	// An event was dispatched.
	Dispatch uint8 = 0
	// Fired periodically by the client to keep the connection alive.
	Heartbeat uint8 = 1
	// Starts a new session during the initial handshake.
	Identify uint8 = 2
	// Update the client's presence.
	PresenceUpdate uint8 = 3
	// Used to join/leave or move between voice channels.
	VoiceStateUpdate uint8 = 4
	// Resume a previous session that was disconnected.
	Resume uint8 = 6
	// You should attempt to reconnect and resume immediately.
	Reconnect uint8 = 7
	// Request information about offline guild members in a large guild.
	RequestGuildMembers uint8 = 8
	// The session has been invalidated. You should reconnect and identify/resume accordingly.
	InvalidSession uint8 = 9
	// Sent immediately after connecting, contains the heartbeat_interval to use.
	Hello uint8 = 10
	// Sent in response to receiving a heartbeat to acknowledge that it has been received.
	HeartbeatAck uint8 = 11
)
