package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var requestCount int32 = 0

func handler(w http.ResponseWriter, r *http.Request) {
	serverName := os.Getenv("SERVER_NAME")

	// Increment request count atomically
	atomic.AddInt32(&requestCount, 1)
	currentCount := atomic.LoadInt32(&requestCount)

	log.Printf("Server %s received a request. Total requests handled: %d\n", serverName, currentCount)
	fmt.Fprintf(w, "Hello from %s. Total requests handled: %d\n", serverName, currentCount)
}

func main() {
	// Ask user for server name and port
	fmt.Print("Enter server name (e.g., Server1): ")
	var serverName string
	fmt.Scanln(&serverName)

	fmt.Print("Enter server port (e.g., 8081): ")
	var port string
	fmt.Scanln(&port)

	// Set environment variables for the server
	os.Setenv("SERVER_NAME", serverName)

	http.HandleFunc("/", handler)
	log.Printf("%s is running on port %s\n", serverName, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
