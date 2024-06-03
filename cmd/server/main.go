package main

import (
	srv "bringyour-test/server"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Get the values of the environment variables
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		log.Println("Failed to load enviroment.")
		return
	}

	// Create a WaitGroup
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	server := srv.Create()
	go server.Run(serverPort)

	// Wait for server closed
	waitGroup.Wait()

	log.Println("Server finished.")
}
