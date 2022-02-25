package main

import (
	"bothoi/gateway"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	gateway.Connect()
}
