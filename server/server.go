package srv

import (
	conn "bringyour-test/pkgs/connection"
	"bringyour-test/pkgs/consts"
	"log"
	"math/rand"
	"net"
	"sync"

	"github.com/google/uuid"
)

type Server struct {
	sync.Map
	ClientUUIDs    sync.Map
	ActiveHandlers sync.Map
}

func Create() *Server {
	return &Server{
		ClientUUIDs:    sync.Map{},
		ActiveHandlers: sync.Map{},
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

var max = 0

func ValueString(values *sync.Map, key any) string {
	data, ok := values.Load(key)
	if !ok {
		return ""
	}
	value, ok := data.(string)
	if !ok {
		return ""
	}
	return value
}

func ValueInt(values *sync.Map, key any) int {
	data, ok := values.Load(key)
	if !ok {
		return 0
	}
	value, ok := data.(int)
	if !ok {
		return 0
	}
	return value
}

func ValueHandler(values *sync.Map, key any) *conn.ConnectionHandler {
	data, ok := values.Load(key)
	if !ok {
		return nil
	}
	value, ok := data.(*conn.ConnectionHandler)
	if !ok {
		return nil
	}
	return value
}

func (server *Server) DeleteClient(clientUUID string) {
	size := ValueInt(&server.ClientUUIDs, "size")
	if size > max {
		max = size
	} else {
		max++
	}
	index := ValueInt(&server.ClientUUIDs, clientUUID)
	lastUUID := ValueString(&server.ClientUUIDs, size-1)
	server.ClientUUIDs.Store(lastUUID, index)
	server.ClientUUIDs.Store(index, lastUUID)
	server.ClientUUIDs.Delete(clientUUID)
	server.ClientUUIDs.Delete(size - 1)
	server.ActiveHandlers.Delete(clientUUID)
	server.ClientUUIDs.Store("size", size-1)
	log.Println(clientUUID, max, "- clients -", size)
}

func (server *Server) HandleConnection(handler *conn.ConnectionHandler) {
	defer handler.Close()

	clientUUID := uuid.New().String()

	size := ValueInt(&server.ClientUUIDs, "size")
	server.ClientUUIDs.Store("size", size+1)
	server.ClientUUIDs.Store(clientUUID, size)
	server.ClientUUIDs.Store(size, clientUUID)
	server.ActiveHandlers.Store(clientUUID, handler)

	// Buffered channel to store received data from the client
	for {
		receivedMessage, err := handler.Receive()
		if err != nil {
			server.DeleteClient(clientUUID)
			return
		}

		if receivedMessage.Prefix == "message" {
			randIndex := rand.Intn(ValueInt(&server.ClientUUIDs, "size"))
			randUUID := ValueString(&server.ClientUUIDs, randIndex)
			randHandler := ValueHandler(&server.ActiveHandlers, randUUID)
			if randHandler != nil {
				randHandler.Send(receivedMessage)
			}
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
				selectedHandler := ValueHandler(&server.ActiveHandlers, uuid)
				if selectedHandler != nil {
					selectedHandler.Send(receivedMessage)
				}
			}
		} else if receivedMessage.Prefix == "close" {
			server.DeleteClient(clientUUID)
			handler.Send(receivedMessage)
			return
		}
	}
}
