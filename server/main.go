package main

import (
	conn "bring-your-test/connection"
	"log"
	"math/rand"
	"net"
	"os"

	"github.com/google/uuid"
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

	// Listen on a TCP port
	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		log.Println("Failed to listen on port:", err)
		return
	}
	defer listener.Close()

	log.Println("Server listening on port " + serverPort)

	clientUUIDs := []string{}
	activeHandlers := map[string]*conn.ConnectionHandler{}

	for {
		// Wait for a handler
		handler, err := conn.CreateByListenr(listener)
		if err != nil {
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(handler, &clientUUIDs, &activeHandlers)
	}
}

var mockUUID = "########-####-####-####-############"

func deleteClient(clientUUID string, clientUUIDs *[]string, activeHandlers *map[string]*conn.ConnectionHandler) {
	log.Println(clientUUID, "- before -", len(*clientUUIDs), len(*activeHandlers))
	newClientUUIDs := []string{}
	for _, client := range *clientUUIDs {
		if client != clientUUID {
			newClientUUIDs = append(newClientUUIDs, client)
		}
	}
	*clientUUIDs = newClientUUIDs
	_, ok := (*activeHandlers)[clientUUID]
	if ok {
		delete((*activeHandlers), clientUUID)
	}
	log.Println(clientUUID, "- after -", len(*clientUUIDs), len(*activeHandlers))
}

func handleConnection(handler *conn.ConnectionHandler, clientUUIDs *[]string, activeHandlers *map[string]*conn.ConnectionHandler) {
	defer handler.Close()

	clientUUID := uuid.New().String()

	(*clientUUIDs) = append((*clientUUIDs), clientUUID)
	(*activeHandlers)[clientUUID] = handler

	// Buffered channel to store received data from the client
	for {
		receivedMessage, err := handler.Receive()
		if err != nil {
			log.Println("- error -", err)
			deleteClient(clientUUID, clientUUIDs, activeHandlers)
			return
		}

		if receivedMessage.Prefix == "message" {
			randUUID := (*clientUUIDs)[rand.Intn(len(*clientUUIDs))]
			if randUUID != clientUUID {
				receivedMessage.UUID = randUUID
			}
			(*activeHandlers)[randUUID].Send(receivedMessage)
		} else if receivedMessage.Prefix == "ok" {
			log.Println("- if prefix == ok -")
			if receivedMessage.UUID == mockUUID {
				log.Println("- if uuid == mockUUID -")
				handler.Send(receivedMessage)
				deleteClient(clientUUID, clientUUIDs, activeHandlers)
				receivedMessage.Prefix = "close"
				handler.Send(receivedMessage)
				return
			} else {
				uuid := receivedMessage.UUID
				receivedMessage.UUID = mockUUID
				(*activeHandlers)[uuid].Send(receivedMessage)
			}
		} else if receivedMessage.Prefix == "close" {
			deleteClient(clientUUID, clientUUIDs, activeHandlers)
			handler.Send(receivedMessage)
			return
		}
	}
}
