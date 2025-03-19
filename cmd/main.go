package main

import (
	"log"

	"github.com/template/go-backend-gin-orm/startup"
)

func main() {
	// Run server
	server := startup.NewServer()
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
