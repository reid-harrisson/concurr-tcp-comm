package main

import (
	cli "bring-your-test/client/client"
	"log"
	"os"
	"strconv"
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
	clientCount, err := strconv.Atoi(os.Getenv("CLIENT_COUNT"))
	if err != nil {
		log.Println("Failed to load enviroment.")
		return
	}
	connectSessionTime, err := strconv.Atoi(os.Getenv("CLIENT_COUNT"))
	if err != nil {
		log.Println("Failed to load enviroment.")
		return
	}

	// Create a WaitGroup
	var waitGroup sync.WaitGroup

	// Add 3 goroutines to the WaitGroup
	log.Println("waitGroup +", clientCount)
	waitGroup.Add(clientCount)

	clients := []*cli.Client{}

	// Create clients
	for i := 0; i < clientCount; i++ {
		clients = append(clients, cli.Create())
	}

	// Run clients
	for _, client := range clients {
		go client.Run(serverPort, &waitGroup, connectSessionTime)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	log.Println("All goroutines finished")
}
