package main

import (
	conn "bring-your-test/connection"
	msg "bring-your-test/models"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var clients = 0
var clientChannel = make(chan int)

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

	// Create a channel
	channel := make(chan string)

	clients = clientCount
	// Run the goroutines
	for i := 0; i < clientCount; i++ {
		go handleConnection(serverPort, &waitGroup, "", channel, connectSessionTime)
	}

	// Retry connections
	for {
		select {
		case message := <-channel:
			waitGroup.Add(1)
			go handleConnection(serverPort, &waitGroup, message, channel, connectSessionTime)
		case count := <-clientChannel:
			log.Println("--------------------------------", count)
			if count == 0 {
				return
			}
		}
	}

	// for message := range channel {
	// 	waitGroup.Add(1)
	// 	go handleConnection(serverPort, &waitGroup, message, channel, connectSessionTime)
	// }

	// Wait for all goroutines to finish
	waitGroup.Wait()

	log.Println("All goroutines finished")
}

var (
	mockMessages = []string{
		"Simplicity is the ultimate sophistication.",
		"The noblest pleasure is the joy of understanding.",
		"Realize that everything connects to everything else.",
		"Learning never exhausts the mind.",
		"The greatest deception men suffer is from their own opinions.",
		"Intellectual passion drives out sensuality.",
		"Small steps are the fastest path to achieving great goals.",
		"It had long since come to my attention that people of accomplishment rarely sat back and let things happen to them. They went out and happened to things.",
		"Obstacles cannot crush me. Every obstacle yields to stern resolve.",
		"I have been impressed with the urgency of doing. Knowing is not enough; we must apply. Being willing is not enough; we must do.",
		"Art is never finished, only abandoned.",
		"The greatest geniuses sometimes accomplish more when they work less.",
		"The function of muscle is to pull, not to push, except in the case of the genitals and the tongue.",
		"The human foot is a masterpiece of engineering and a work of art.",
		"I love those who can smile in trouble, who can gather strength from distress, and grow brave by reflection.",
		"Time stays long enough for those who use it.",
		"Painting is poetry that is seen rather than felt, and poetry is painting that is felt rather than seen.",
		"The human bird shall take his first flight, filling the world with amazement, all writings with his fame, and bringing eternal glory to the nest whence he sprang.",
		"Dwell on the beauty of life. Watch the stars, and see yourself running with them.",
		"You can have no dominion greater or less than that over yourself.",
		"The most beautiful thing we can experience is the mysterious. It is the source of all true art and all science.",
		"Imagination is more important than knowledge.",
		"The only source of knowledge is experience.",
		"The only way to do great work is to love what you do.",
		"Everything should be made as simple as possible, but not simpler.",
		"The important thing is not to stop questioning. Curiosity has its own reason for existing.",
		"The true sign of intelligence is not knowledge but imagination.",
		"In the middle of difficulty lies opportunity.",
		"The most beautiful experience we can have is the mysterious.",
		"I have no special talent. I am only passionately curious.",
		"The only limit to our realization of tomorrow will be our doubts of today.",
		"The world as we have created it is a process of our thinking. It cannot be changed without changing our thinking.",
		"Try not to become a man of success, but rather try to become a man of value.",
		"The measure of intelligence is the ability to change.",
		"Education is what remains after one has forgotten what one has learned in school.",
		"The person who follows the crowd will usually go no further than the crowd. The person who walks alone is likely to find himself in places no one has ever been.",
		"Two things are infinite: the universe and human stupidity; and I'm not sure about the universe.",
		"Insanity is doing the same thing over and over again and expecting different results.",
		"Look deep into nature, and then you will understand everything better.",
		"The most incomprehensible thing about the universe is that it is comprehensible.",
	}
	mockUUID = "########-####-####-####-############"
)

func closeClient(waitGroup *sync.WaitGroup, handler *conn.ConnectionHandler) {
	handler.Close()
	waitGroup.Done()
	log.Println("--- close client ---")
}

func handleConnection(serverPort string, waitGroup *sync.WaitGroup, message string, reconnectChannel chan string, connectSessionTime int) {
	// Connect to the TCP server
	handler, err := conn.Create(":" + serverPort)
	if err != nil {
		log.Println("--- failed to create ---")
		waitGroup.Done()
		return
	}

	// Initialize message
	if message == "" {
		message = mockMessages[rand.Intn(len(mockMessages))]
	}

	// Close connection after 60s
	timer := time.NewTimer(time.Duration(connectSessionTime * int(time.Second)))
	defer timer.Stop()

	// Buffered channel to store received data from the client
	messageChannel := make(chan msg.Message, 1)

	err = handler.Send(msg.Message{
		UUID:    mockUUID,
		Prefix:  "message",
		Message: message,
	})
	if err != nil {
		closeClient(waitGroup, handler)
		return
	}

	go func(handler *conn.ConnectionHandler, channel chan msg.Message, waitGroup *sync.WaitGroup) {
		for {
			receivedMessage, err := handler.Receive()
			if err != nil {
				log.Println("-- failed to receive --", err)
				break
			}

			if receivedMessage.Prefix == "close" {
				break
			}
			messageChannel <- receivedMessage
			if receivedMessage.Prefix == "ok" && receivedMessage.UUID == mockUUID {
				break
			}
		}
		close(messageChannel)
		log.Println("-- receive goroutine finished --")
	}(handler, messageChannel, waitGroup)

	for {
		select {
		case receivedMessage, isOpened := <-messageChannel:
			if !isOpened {
				log.Println("--- message channel closed ---")
				waitGroup.Done()
				return
			}
			if receivedMessage.Prefix == "ok" && receivedMessage.UUID == mockUUID {
				log.Println("--- close by OX condition ---")
				closeClient(waitGroup, handler)
				clients--
				clientChannel <- clients
				return
			}
			if receivedMessage.Prefix == "message" {
				receivedMessage.Prefix = "ok"
			}
			handler.Send(receivedMessage)

		case <-timer.C:
			log.Println("--- close by timer ---")
			handler.Send(msg.Message{
				Prefix:  "close",
				UUID:    mockUUID,
				Message: "",
			})
			closeClient(waitGroup, handler)
			reconnectChannel <- message
			return
		}
	}
}
