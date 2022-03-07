package app_command_option_type

type AppCommandOptionType byte

const (
	SubCommand AppCommandOptionType = iota + 1
	SubCommandGroup
	String
	Integer
	Boolean
	User
	Channel
	Role
	Mentionable
	Number
	Attachment
)
