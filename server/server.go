package srv

import (
	conn "bring-your-test/pkgs/connection"
	"bring-your-test/pkgs/consts"
	"log"
	"math/rand"
	"net"

	"github.com/google/uuid"
)

type Server struct {
	ClientUUIDs    []string
	ActiveHandlers map[string]*conn.ConnectionHandler
}

func Create() *Server {
	return &Server{
		ClientUUIDs:    []string{},
		ActiveHandlers: map[string]*conn.ConnectionHandler{},
	}
}

func (server *Server) Run(serverPort string) {
	// Listen on a TCP port
	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		log.Println("Failed to listen on port:", err)
		return
	}
	defer listener.Close()

	log.Println("Server listening on port " + serverPort)

	for {
		// Wait for a handler
		handler, err := conn.CreateByListenr(listener)
		if err != nil {
			continue
		}

		// Handle the connection in a new goroutine
		go server.HandleConnection(handler)
	}
}

func (server *Server) DeleteClient(clientUUID string) {
	// log.Println(clientUUID, "- before -", len(server.ClientUUIDs), len(server.ActiveHandlers))
	newClientUUIDs := []string{}
	for _, client := range server.ClientUUIDs {
		if client != clientUUID {
			newClientUUIDs = append(newClientUUIDs, client)
		}
	}
	server.ClientUUIDs = newClientUUIDs
	_, ok := server.ActiveHandlers[clientUUID]
	if ok {
		delete(server.ActiveHandlers, clientUUID)
	}
	// log.Println(clientUUID, "- after -", len(server.ClientUUIDs), len(server.ActiveHandlers))
}

func (server *Server) HandleConnection(handler *conn.ConnectionHandler) {
	defer handler.Close()

	clientUUID := uuid.New().String()

	server.ClientUUIDs = append(server.ClientUUIDs, clientUUID)
	server.ActiveHandlers[clientUUID] = handler

	// Buffered channel to store received data from the client
	for {
		receivedMessage, err := handler.Receive()
		if err != nil {
			server.DeleteClient(clientUUID)
			return
		}

		if receivedMessage.Prefix == "message" {
			randUUID := server.ClientUUIDs[rand.Intn(len(server.ClientUUIDs))]
			if randUUID != clientUUID {
				receivedMessage.UUID = randUUID
			}
			server.ActiveHandlers[randUUID].Send(receivedMessage)
		} else if receivedMessage.Prefix == "ok" {
			if receivedMessage.UUID == consts.MockUUID {
				handler.Send(receivedMessage)
				server.DeleteClient(clientUUID)
				receivedMessage.Prefix = "close"
				handler.Send(receivedMessage)
				return
			} else {
				uuid := receivedMessage.UUID
				receivedMessage.UUID = consts.MockUUID
				server.ActiveHandlers[uuid].Send(receivedMessage)
			}
		} else if receivedMessage.Prefix == "close" {
			server.DeleteClient(clientUUID)
			handler.Send(receivedMessage)
			return
		}
	}
}
