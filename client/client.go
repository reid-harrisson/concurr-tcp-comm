package cli

import (
	conn "bringyour-test/pkgs/connection"
	consts "bringyour-test/pkgs/consts"
	msg "bringyour-test/pkgs/models"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Client struct {
	Index            int
	ReconnectChannel chan string
	MessagesReceived []msg.Message
	MessagesSent     []msg.Message
}

func Create(index int) *Client {
	return &Client{
		Index:            index,
		MessagesReceived: []msg.Message{},
		MessagesSent:     []msg.Message{},
		ReconnectChannel: make(chan string),
	}
}

func (client *Client) HandleConnection(serverPort string, message string, sessionTime int, delayTime int) {
	// Connect to the TCP server
	handler, err := conn.Create(":" + serverPort)
	if err != nil {
		log.Printf("Client %d failed to connect server.\n", client.Index)
		client.ReconnectChannel <- ""
		return
	}
	log.Printf("Client %d connected to server.", client.Index)

	// Close connection after 60s
	timer := time.NewTimer(time.Second * time.Duration(sessionTime))
	defer timer.Stop()

	// Buffered channel to store received data from the client
	messageChannel := make(chan msg.Message, 1)

	err = handler.Send(msg.Message{
		UUID:    consts.MockUUID,
		Prefix:  "message",
		Message: message,
	})
	log.Printf("Client %d sent messge X.", client.Index)
	if err != nil {
		handler.Close()
		client.ReconnectChannel <- ""
		return
	}

	go func() {
		for {
			receivedMessage, err := handler.Receive()
			if err != nil {
				break
			}

			if receivedMessage.Prefix == "close" {
				break
			}

			// Update message received history
			client.MessagesReceived = append(client.MessagesReceived, receivedMessage)

			messageChannel <- receivedMessage

			// Close routine in OX condition
			if receivedMessage.Prefix == "ok" && receivedMessage.UUID == consts.MockUUID {
				break
			}
		}
		close(messageChannel)
		// log.Println("-- receive goroutine finished --")
	}()

	for {
		select {
		case receivedMessage := <-messageChannel:
			log.Printf("Client %d received %s.\n", client.Index, consts.ShortMessage(receivedMessage))
			if receivedMessage.Prefix == "ok" && receivedMessage.UUID == consts.MockUUID {
				log.Printf("Client %d successfully finished.\n", client.Index)
				handler.Close()
				client.ReconnectChannel <- ""
				return
			}
			if receivedMessage.Prefix == "message" {
				receivedMessage.Prefix = "ok"
			}
			if receivedMessage.Prefix == "ok" {
				randTime := rand.Intn(delayTime) + 1
				time.Sleep(time.Second * time.Duration(randTime))
				handler.Send(receivedMessage)
			}
			log.Printf("Client %d sent %s.\n", client.Index, consts.ShortMessage(receivedMessage))
			client.MessagesSent = append(client.MessagesSent, receivedMessage)

		case <-timer.C:
			handler.Send(msg.Message{
				Prefix:  "close",
				UUID:    consts.MockUUID,
				Message: "",
			})
			handler.Close()
			log.Printf("Client %d's connection attempt timed out.\n", client.Index)
			client.ReconnectChannel <- message
			return
		}
	}
}

func (client *Client) Run(serverPort string, waitGroup *sync.WaitGroup, sessionTime int, delayTime int) {
	go client.HandleConnection(serverPort, consts.RandomMessage(), sessionTime, delayTime)
	for message := range client.ReconnectChannel {
		if message == "" {
			break
		}
		go client.HandleConnection(serverPort, message, sessionTime, delayTime)
	}
	waitGroup.Done()
}
