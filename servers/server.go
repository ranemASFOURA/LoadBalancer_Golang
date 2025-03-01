package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"load-balancer-project/config"
)

var requestCount int
var requestMu sync.Mutex // Mutex to protect shared variable access

func main() {
	// Ensure a server name is provided as an argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run server.go <ServerName>")
	}

	serverName := os.Args[1]

	// Load configuration from config.json
	configData, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var serverURL string
	found := false
	for _, server := range configData.Servers {
		if server.Name == serverName {
			serverURL = server.URL
			found = true
			break
		}
	}

	// Exit if the server name is not found in the config file
	if !found {
		log.Fatalf("Server %s not found in config.json", serverName)
	}

	// Parse the server URL to extract the port number
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		log.Fatalf("Invalid URL for server %s: %v", serverName, err)
	}

	port := parsedURL.Port()
	if port == "" {
		log.Fatalf("No port found in URL for server %s", serverName)
	}

	// Health check endpoint to verify server availability
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s received a HealthCheck request.\n", serverName)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HealthCheck OK"))
	})

	// Main request handler with simulated variable processing time
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Increment request count safely using a mutex
		requestMu.Lock()
		requestCount++
		fmt.Printf("%s received a request. Total requests handled: %d\n", serverName, requestCount)
		requestMu.Unlock()

		// Simulate different processing times for each request
		processTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
		time.Sleep(processTime)

		// Respond with processing time and server name
		fmt.Fprintf(w, "Hello from %s (Processing time: %v)\n", serverName, processTime)
	})

	// Start the HTTP server
	log.Printf("%s is running on port %s\n", serverName, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
