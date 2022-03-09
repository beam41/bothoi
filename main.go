package main

import (
	"bothoi/commands"
	"bothoi/gateway"
)

func main() {
	commands.Register()
	gateway.Connect()
}
