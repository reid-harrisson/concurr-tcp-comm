package main

import (
	cli "bringyour-test/client"
	"bringyour-test/pkgs/consts"
	srv "bringyour-test/server"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {

	// Get the values of the environment variables
	serverPort := "8000"
	assert.Equal(t, "8000", serverPort)

	clientCount, err := strconv.Atoi("1")
	if assert.Equal(t, nil, err) {
		assert.Equal(t, 1, clientCount)
	}
	if err != nil {
		log.Println("Failed to load enviroment.")
		return
	}
	sessionTime := 10
	delayTime := 5

	server := srv.Create()
	go server.Run(serverPort)

	time.Sleep(time.Second * 2)

	// Create a WaitGroup
	var waitGroup sync.WaitGroup
	waitGroup.Add(clientCount)

	// Initialize
	clients := []*cli.Client{}

	// Create clients
	for i := 0; i < clientCount; i++ {
		clients = append(clients, cli.Create(i+1))
	}

	// Run clients
	for _, client := range clients {
		go client.Run(serverPort, &waitGroup, sessionTime, delayTime)
	}

	// Wait for all clients finished
	waitGroup.Wait()
	log.Println("All clients finished.")

	assert.Equal(t, consts.MockUUID, clients[0].MessagesReceived[0].UUID)
	assert.Equal(t, "message", clients[0].MessagesReceived[0].Prefix)
	assert.Equal(t, consts.MockUUID, clients[0].MessagesReceived[1].UUID)
	assert.Equal(t, "ok", clients[0].MessagesReceived[1].Prefix)

	// Wait for server closed
	log.Println("Server finished.")
}
