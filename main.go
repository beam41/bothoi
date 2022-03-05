package main

import (
	"bothoi/app_command"
	"bothoi/gateway"
)

func main() {
	app_command.Register()
	gateway.Connect()
}
