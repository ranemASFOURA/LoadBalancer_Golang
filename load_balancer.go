package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Server struct represents a backend server's configuration
type Server struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Config struct holds load balancer settings
type Config struct {
	Servers             []Server `yaml:"servers"`
	HealthCheckInterval int      `yaml:"health_check_interval"`
}

// ServerState keeps track of each server's status
type ServerState struct {
	Name        string
	URL         string
	Connections int
	Alive       bool
}

// LoadBalancer struct manages server selection and load distribution
type LoadBalancer struct {
	servers []*ServerState
	mu      sync.Mutex
}

// NewLoadBalancer initializes a Load Balancer with given servers
func NewLoadBalancer(servers []Server) *LoadBalancer {
	lb := &LoadBalancer{}
	for _, srv := range servers {
		lb.servers = append(lb.servers, &ServerState{
			Name:  srv.Name,
			URL:   srv.URL,
			Alive: true,
		})
	}
	return lb
}

// LoadConfig reads the configuration from a YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetLeastConnectionsServer selects the backend server with the least connections
func (lb *LoadBalancer) GetLeastConnectionsServer() *ServerState {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	var bestServer *ServerState
	for _, server := range lb.servers {
		if server.Alive && (bestServer == nil || server.Connections < bestServer.Connections) {
			bestServer = server
		}
	}
	if bestServer != nil {
		bestServer.Connections++
		fmt.Printf("Load Balancer: Redirecting request to %s (Active connections: %d)\n", bestServer.Name, bestServer.Connections)
	}
	return bestServer
}

// ReleaseConnection decrements the connection count after a request is completed
func (lb *LoadBalancer) ReleaseConnection(server *ServerState) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	server.Connections--
}

// HealthCheck continuously checks the health of backend servers
func (lb *LoadBalancer) HealthCheck(interval int) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		lb.mu.Lock()
		for _, server := range lb.servers {
			resp, err := http.Get(server.URL)
			server.Alive = (err == nil && resp.StatusCode == http.StatusOK)
		}
		lb.mu.Unlock()
	}
}

// ServeHTTP forwards requests to the appropriate backend server
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.GetLeastConnectionsServer()
	if server == nil {
		http.Error(w, "All servers are unavailable", http.StatusServiceUnavailable)
		return
	}

	// Handle requests concurrently
	go func(srv *ServerState) {
		defer lb.ReleaseConnection(srv)

		resp, err := http.Get(srv.URL)
		if err != nil {
			http.Error(w, "Error connecting to the server", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Display the server name in the response
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintf(w, "Your request was routed to: %s\n", srv.Name)
	}(server)
}

func main() {
	// Load configuration from file
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Ask user to enter the Load Balancer port
	fmt.Print("Enter the port for the Load Balancer: ")
	var portStr string
	fmt.Scanln(&portStr)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Invalid port number")
	}

	// Start Load Balancer
	lb := NewLoadBalancer(config.Servers)
	go lb.HealthCheck(config.HealthCheckInterval)

	http.Handle("/", lb)
	address := fmt.Sprintf(":%d", port)
	log.Printf("Load Balancer is running on port %d\n", port)
	log.Fatal(http.ListenAndServe(address, nil))
}
