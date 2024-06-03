package srv

import (
	conn "bringyour-test/pkgs/connection"
	"bringyour-test/pkgs/consts"
	"log"
	"math/rand"
	"net"

	"github.com/google/uuid"
)

type Server struct {
	ClientUUIDs    map[string]int
	ClientIndex    map[int]string
	ActiveHandlers map[string]*conn.ConnectionHandler
}

func Create() *Server {
	return &Server{
		ClientIndex:    map[int]string{},
		ClientUUIDs:    map[string]int{},
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
	log.Println(clientUUID, "- before -", len(server.ClientUUIDs), len(server.ClientIndex), len(server.ActiveHandlers))
	size := len(server.ClientUUIDs)
	index := server.ClientUUIDs[clientUUID]
	lastUUID := server.ClientIndex[size-1]
	server.ClientUUIDs[lastUUID] = index
	server.ClientIndex[index] = lastUUID
	delete(server.ClientUUIDs, clientUUID)
	delete(server.ClientIndex, size-1)
	delete(server.ActiveHandlers, clientUUID)
	log.Println(clientUUID, "- after -", len(server.ClientUUIDs), len(server.ClientIndex), len(server.ActiveHandlers))
}

func (server *Server) HandleConnection(handler *conn.ConnectionHandler) {
	defer handler.Close()

	clientUUID := uuid.New().String()

	server.ClientUUIDs[clientUUID] = len(server.ClientUUIDs)
	server.ClientIndex[len(server.ClientIndex)] = clientUUID
	server.ActiveHandlers[clientUUID] = handler

	// Buffered channel to store received data from the client
	for {
		receivedMessage, err := handler.Receive()
		if err != nil {
			server.DeleteClient(clientUUID)
			return
		}

		if receivedMessage.Prefix == "message" {
			randIndex := rand.Intn(len(server.ClientUUIDs))
			randUUID := server.ClientIndex[randIndex]
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
