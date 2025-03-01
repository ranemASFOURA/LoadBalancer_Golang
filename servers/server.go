package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"load-balancer-project/config"
)

// Global variables to keep track of the number of requests
var requestCount int
var requestMu sync.Mutex // Mutex to protect shared resource (requestCount)

func main() {
	// Ensure a server name is provided as a command-line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run server.go <ServerName>")
	}

	// Retrieve the server name from the command-line arguments
	serverName := os.Args[1]

	// Load the server configuration from the config file
	configData, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Search for the specified server in the configuration
	var serverURL string
	found := false
	for _, server := range configData.Servers {
		if server.Name == serverName {
			serverURL = server.URL
			found = true
			break
		}
	}

	// If the server is not found in the config file, exit with an error
	if !found {
		log.Fatalf("Server %s not found in config.json", serverName)
	}

	// Extract the port from the server URL
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		log.Fatalf("Invalid URL for server %s: %v", serverName, err)
	}

	// Get the port from the URL. If no port is found, log an error and exit.
	port := parsedURL.Port()
	if port == "" {
		log.Fatalf("No port found in URL for server %s", serverName)
	}

	// Define the health check route
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		// Log the health check request
		fmt.Printf("%s received a HealthCheck request.\n", serverName)
		// Respond with a successful health check status
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HealthCheck OK"))
	})

	// Define the main route to handle other requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Lock the mutex to safely update the shared request count
		requestMu.Lock()
		requestCount++
		// Log the request received and the total request count
		fmt.Printf("%s received a request. Total requests handled: %d\n", serverName, requestCount)
		// Unlock the mutex after updating the request count
		requestMu.Unlock()
		// Respond with a message indicating the server handling the request
		fmt.Fprintf(w, "Hello from %s\n", serverName)
	})

	// Log that the server is running and listening on the specified port
	log.Printf("%s is running on port %s\n", serverName, port)
	// Start the HTTP server and listen for incoming requests
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
