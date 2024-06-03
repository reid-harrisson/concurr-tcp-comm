package conn

import (
	"bring-your-test/pkgs/consts"
	msg "bring-your-test/pkgs/models"
	"encoding/json"
	"io"
	"log"
	"net"
	"strings"
)

type ConnectionHandler struct {
	Connection net.Conn
}

func ShortMessage(message msg.Message) string {
	if message.Prefix == "message" && message.UUID == consts.MockUUID {
		return "MX"
	} else if message.Prefix == "message" && message.UUID != consts.MockUUID {
		return "MY"
	} else if message.Prefix == "ok" && message.UUID == consts.MockUUID {
		return "OX"
	} else if message.Prefix == "ok" && message.UUID != consts.MockUUID {
		return "OY"
	}
	return strings.ToUpper(message.Prefix)
}

func Create(url string) (*ConnectionHandler, error) {
	connection, err := net.Dial("tcp", url)
	if err != nil {
		log.Println("Failed to connect to server:", err)
		return nil, err
	}
	log.Println("-", connection.LocalAddr(), connection.RemoteAddr(), "- create -")
	return &ConnectionHandler{
		Connection: connection,
	}, nil
}

func CreateByListenr(listener net.Listener) (*ConnectionHandler, error) {
	connection, err := listener.Accept()
	if err != nil {
		log.Println("Failed to accept connection:", err)
		return nil, err
	}
	log.Println("-", connection.LocalAddr(), connection.RemoteAddr(), "- create -")
	return &ConnectionHandler{
		Connection: connection,
	}, nil
}

func (handler *ConnectionHandler) Close() error {
	log.Println("-", handler.Connection.LocalAddr(), handler.Connection.RemoteAddr(), "- close -")
	return handler.Connection.Close()
}

func (handler *ConnectionHandler) Send(modelMessage msg.Message) error {
	// Create modelMessage
	jsonMessage, _ := json.Marshal(modelMessage)

	// Send data to the server
	_, err := handler.Connection.Write([]byte(jsonMessage))
	if err != nil {
		log.Println("Error writing to connection:", err)
		return err
	}
	log.Println("-", handler.Connection.LocalAddr(), handler.Connection.RemoteAddr(), "- sent -", ShortMessage(modelMessage), "-")
	return nil
}

func (handler *ConnectionHandler) Receive() (msg.Message, error) {
	// Read the response from the server
	response := make([]byte, 1024)
	size, err := handler.Connection.Read(response)
	if size == 0 || err == io.EOF {
		// log.Println("-- disconnected")
	}
	if err != nil {
		return msg.Message{}, err
	}
	receivedMessage := msg.Message{}
	json.Unmarshal(response[:size], &receivedMessage)
	log.Println("-", handler.Connection.LocalAddr(), handler.Connection.RemoteAddr(), "- receive -", ShortMessage(receivedMessage), "-")
	return receivedMessage, nil
}
