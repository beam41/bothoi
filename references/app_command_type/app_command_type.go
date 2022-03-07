package app_command_type

type AppCommandType byte

const (
	ChatInput AppCommandType = iota + 1
	User
	Message
)
