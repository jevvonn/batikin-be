package main

import (
	"batikin-be/internal/bootstrap"
	"log"
)

func main() {
	err := bootstrap.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
