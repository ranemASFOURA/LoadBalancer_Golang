package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// simulateRequests function to simulate multiple requests to the load balancer
func simulateRequests(numRequests int) {
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			resp, err := http.Get("http://localhost:8080") // Load Balancer address
			if err != nil {
				log.Printf("Error sending request %d: %v", requestID, err)
				return
			}
			defer resp.Body.Close()

			// Log the response status
			fmt.Printf("Request %d received response status: %s\n", requestID, resp.Status)
		}(i)
	}

	// Wait for all requests to finish
	wg.Wait()
}

func main() {
	// Simulate sending 10 requests to the Load Balancer
	simulateRequests(5)
}
