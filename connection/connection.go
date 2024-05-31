package conn

import (
	msg "bring-your-test/models"
	"encoding/json"
	"log"
	"net"
)

type ConnectionHandler struct {
	Connection net.Conn
}

func Create(url string) (*ConnectionHandler, error) {
	connection, err := net.Dial("tcp", url)
	if err != nil {
		log.Println("Failed to connect to server:", err)
		return nil, err
	}
	return &ConnectionHandler{
		Connection: connection,
	}, nil
}

func (handler *ConnectionHandler) Close() {
	handler.Connection.Close()
}

func (handler *ConnectionHandler) Send(uuid string, prefix string, message string) error {
	// Create modelMessage
	modelMessage := msg.Message{
		UUID:   uuid,
		Prefix: prefix,
		Data:   message,
	}
	jsonMessage, _ := json.Marshal(modelMessage)

	// Send data to the server
	_, err := handler.Connection.Write([]byte(jsonMessage))
	if err != nil {
		log.Println("Error writing to connection:", err)
		return err
	}
	return nil
}

func (handler *ConnectionHandler) Receive() (msg.Message, error) {
	// Read the response from the server
	response := make([]byte, 1024)
	size, err := handler.Connection.Read(response)
	if err != nil {
		return msg.Message{}, err
	}
	modelMessage := msg.Message{}
	err = json.Unmarshal(response[:size], &modelMessage)
	if err != nil {
		return msg.Message{}, err
	}
	return modelMessage, nil
}
