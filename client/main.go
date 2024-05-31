package main

import (
	conn "bring-your-test/connection"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

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
	clientCount, err := strconv.Atoi(os.Getenv("CLIENT_COUNT"))
	if err != nil {
		log.Println("Failed to load enviroment.")
		return
	}

	// Create a WaitGroup
	var waitGroup sync.WaitGroup

	// Add 3 goroutines to the WaitGroup
	waitGroup.Add(clientCount)

	// Create a channel
	channel := make(chan string)

	// Run the goroutines
	for i := 0; i < clientCount; i++ {
		go handleConnection(serverPort, &waitGroup, "", channel)
	}

	// Retry connections
	for uuid := range channel {
		waitGroup.Add(1)
		go handleConnection(serverPort, &waitGroup, uuid, channel)
	}

	// Wait for all goroutines to finish
	waitGroup.Wait()

	log.Println("All goroutines finished")
}

var mockMessages = []string{
	"The lion roared loudly in the jungle.",
	"The tiger stealthily stalked its prey.",
	"The elephant sprayed itself with water using its trunk.",
	"The giraffe gracefully walked across the savanna.",
	"The zebra's black and white stripes stood out against the grass.",
	"The monkey swung from branch to branch in the canopy.",
	"The koala slept peacefully in the eucalyptus tree.",
	"The panda chomped on bamboo shoots in the bamboo forest.",
	"The kangaroo hopped across the Australian outback.",
	"The gorilla beat its chest in a display of strength.",
	"The cheetah sprinted across the plains at top speed.",
	"The grizzly bear caught salmon in the rushing river.",
	"The hippopotamus wallowed in the muddy river.",
	"The rhinoceros charged through the grasslands.",
	"The crocodile sunbathed on the riverbank with its jaws wide open.",
	"The alligator's eyes glinted in the water as it lurked.",
	"The leopard camouflaged itself among the tree branches.",
	"The orangutan swung gracefully through the rainforest trees.",
	"The dolphin leaped out of the water in a graceful arc.",
	"The penguin waddled clumsily across the snowy ice.",
	"The octopus changed color to blend in with its surroundings.",
	"The seahorse floated gently among the seagrass.",
	"The jellyfish glowed in the dark depths of the ocean.",
	"The chimpanzee groomed its companion in the treetops.",
	"The ostrich raced across the African savanna.",
	"The flamingo stood on one leg in the shallow waters of the lake.",
	"The polar bear hunted seals on the sea ice.",
	"The red panda curled up in a ball in its tree nest.",
	"The sloth moved slowly through the treetops.",
	"The Komodo dragon basked in the sun on Komodo Island.",
	"The parrot squawked loudly in the jungle canopy.",
	"The toucan's brightly colored beak caught the sunlight.",
	"The platypus swam gracefully in the river.",
	"The meerkat stood guard on its hind legs in the desert.",
	"The armadillo rolled up into a tight ball for protection.",
	"The corgi herded sheep on the farm with enthusiasm.",
	"The tabby cat purred contentedly in its owner's lap.",
	"The blue whale breached the surface of the ocean with a mighty splash.",
	"The wolf howled at the full moon in the night.",
	"The fox darted among the bushes in search of prey.",
	"The otter slid into the water and began to play.",
	"The hedgehog curled into a spiky ball for defense.",
	"The peacock fanned out its iridescent feathers in a display.",
	"The bald eagle soared high in the sky on majestic wings.",
	"The raccoon scavenged for food in the trash cans at night.",
	"The walrus basked on the beach with its tusks gleaming.",
	"The alpaca nibbled on grass in the fields of the Andes.",
	"The llama spit in annoyance at a pesky visitor.",
	"The emu sprinted across the Australian outback.",
	"The pufferfish inflated itself to scare away predators.",
}

func handleConnection(serverPort string, waitGroup *sync.WaitGroup, prevUUID string, channel chan string) {
	// Get the current time
	startTime := time.Now()

	// Generate a new UUID
	clientUUID := prevUUID
	if clientUUID == "" {
		clientUUID = uuid.New().String()
	}

	// Handle goroutine completion
	defer waitGroup.Done()
	defer log.Printf("Client (%s) goroutine finished\n", clientUUID)
	log.Printf("Client (%s) goroutine started\n", clientUUID)

	// Connect to the TCP server
	handler, err := conn.Create(":" + serverPort)
	if err != nil {
		return
	}
	defer handler.Close()

	// Create modelMessage
	handler.Send(clientUUID, "connect", mockMessages[rand.Intn(len(mockMessages))])

	for {
		// Close connection after 60s
		currentTime := time.Now()
		if currentTime.Sub(startTime) >= 5*time.Second {
			defer func(channel chan string, uuid string) {
				channel <- uuid
			}(channel, clientUUID)
			return
		}

		// Read the response from the server
		modelMessage, err := handler.Receive()
		if err == nil {
			log.Println(modelMessage)
		}
	}
}
