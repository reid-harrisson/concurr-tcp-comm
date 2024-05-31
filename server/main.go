package main

import (
	"log"
	"net"
	"os"

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

	for {
		// Wait for a connection
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(connection)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read data from the connection
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Failed to read from connection:", err)
		return
	}

	// Print the received data
	log.Println("Received data:", string(buf[:n]))

	// Write a response to the connection
	_, err = conn.Write([]byte("Hello, client!"))
	if err != nil {
		log.Println("Failed to write to connection:", err)
		return
	}

	// Write a response to the connection
	_, err = conn.Write([]byte("Hello, client!"))
	if err != nil {
		log.Println("Failed to write to connection:", err)
		return
	}
}
