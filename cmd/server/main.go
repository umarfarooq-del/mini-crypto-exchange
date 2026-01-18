package main

import (
	"mini-crypto-exchange/internal/server"

	"log"
)

func main() {
	err := server.RunServer()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}