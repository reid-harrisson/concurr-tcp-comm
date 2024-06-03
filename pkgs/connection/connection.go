package conn

import (
	msg "bringyour-test/pkgs/models"
	"encoding/json"
	"net"
)

type ConnectionHandler struct {
	Connection net.Conn
}

func Create(url string) (*ConnectionHandler, error) {
	connection, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	return &ConnectionHandler{
		Connection: connection,
	}, nil
}

func CreateByListenr(listener net.Listener) (*ConnectionHandler, error) {
	connection, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	return &ConnectionHandler{
		Connection: connection,
	}, nil
}

func (handler *ConnectionHandler) Close() error {
	return handler.Connection.Close()
}

func (handler *ConnectionHandler) Send(modelMessage msg.Message) error {
	// Create modelMessage
	jsonMessage, _ := json.Marshal(modelMessage)

	// Send data to the server
	_, err := handler.Connection.Write([]byte(jsonMessage))
	if err != nil {
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
	receivedMessage := msg.Message{}
	json.Unmarshal(response[:size], &receivedMessage)
	return receivedMessage, nil
}
