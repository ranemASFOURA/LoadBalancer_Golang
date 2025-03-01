package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"load-balancer-project/config"
)

// LoadBalancer struct to manage multiple servers
type LoadBalancer struct {
	servers []config.Server // List of servers
	mu      sync.Mutex      // Mutex to protect shared data
	logger  *log.Logger     // Logger for logging activities
}

// Function to create a new Load Balancer
func NewLoadBalancer(servers []config.Server) *LoadBalancer {
	// Open the log file in append mode, create if doesn't exist
	file, err := os.OpenFile("logfile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	// Create a new logger
	logger := log.New(file, "LoadBalancer: ", log.LstdFlags)

	return &LoadBalancer{
		servers: servers,
		logger:  logger,
	}
}

// Function to get the server with the least active connections
func (lb *LoadBalancer) getLeastConnectionsServer() *config.Server {
	lb.mu.Lock()         // Locking to prevent race conditions
	defer lb.mu.Unlock() // Unlock after function completes

	var leastLoadedServer *config.Server
	minConnections := int(^uint(0) >> 1) // Max int value

	// Loop through servers to find the one with the least active connections
	for i := range lb.servers {
		if lb.servers[i].Healthy && lb.servers[i].ActiveConnections < minConnections {
			leastLoadedServer = &lb.servers[i]
			minConnections = lb.servers[i].ActiveConnections
		}
	}

	return leastLoadedServer
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the server with the least number of active connections
	server := lb.getLeastConnectionsServer()
	if server == nil {
		http.Error(w, "No healthy servers available", http.StatusServiceUnavailable)
		return
	}

	// Increase active connection count for the selected server
	lb.mu.Lock()
	server.ActiveConnections++
	lb.logger.Printf("Request directed to %s - Active Connections: %d\n", server.Name, server.ActiveConnections)
	lb.mu.Unlock()

	// Redirect request to the selected server
	targetURL, err := url.Parse(server.URL)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)

	// Decrease active connection count after request is served
	lb.mu.Lock()
	server.ActiveConnections--
	lb.logger.Printf("Request completed on %s - Active Connections: %d\n", server.Name, server.ActiveConnections)
	lb.mu.Unlock()
}

// Function to check server health at intervals
func (lb *LoadBalancer) HealthCheck(interval time.Duration) {
	for {
		// Check the health of each server in the list
		for i := range lb.servers {
			// Send a health check request to the server
			resp, err := http.Get(lb.servers[i].URL + "/healthcheck")
			if err != nil || resp.StatusCode != http.StatusOK {
				// If server is down, mark as unhealthy
				lb.logger.Printf("Server %s is DOWN\n", lb.servers[i].Name)
				lb.servers[i].Healthy = false
			} else {
				// If server is up, mark as healthy
				lb.servers[i].Healthy = true
				lb.logger.Printf("Server %s is UP\n", lb.servers[i].Name)
			}
		}
		// Wait for the specified interval before checking again
		time.Sleep(interval)
	}
}

func main() {
	// Load configuration from the config file
	configData, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Convert HealthCheckInterval string to time.Duration
	interval, err := time.ParseDuration(configData.HealthCheckInterval)
	if err != nil {
		log.Fatalf("Error parsing HealthCheckInterval: %v", err)
	}

	// Create the Load Balancer with the loaded servers
	loadBalancer := NewLoadBalancer(configData.Servers)

	// Start health checks in a separate Goroutine
	go loadBalancer.HealthCheck(interval)

	// Start the Load Balancer server and listen on the specified port
	log.Printf("Load Balancer is running on port %s\n", configData.ListenPort)
	log.Fatal(http.ListenAndServe(configData.ListenPort, loadBalancer))
}
