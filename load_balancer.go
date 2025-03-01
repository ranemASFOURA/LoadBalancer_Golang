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
	servers []config.Server
	mu      sync.Mutex
	logger  *log.Logger
}

// Function to create a new Load Balancer
func NewLoadBalancer(servers []config.Server) *LoadBalancer {
	// فتح ملف اللوج
	file, err := os.OpenFile("logfile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	logger := log.New(file, "LoadBalancer: ", log.LstdFlags)

	return &LoadBalancer{
		servers: servers,
		logger:  logger,
	}
}

// Function to get the server with the least active connections
func (lb *LoadBalancer) getLeastConnectionsServer() *config.Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	var leastLoadedServer *config.Server
	minConnections := int(^uint(0) >> 1) // Max int value

	for i := range lb.servers {
		if lb.servers[i].Healthy && lb.servers[i].ActiveConnections < minConnections {
			leastLoadedServer = &lb.servers[i]
			minConnections = lb.servers[i].ActiveConnections
		}
	}

	return leastLoadedServer
}

// Function to handle incoming requests and forward them to available servers
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// الحصول على الخادم الذي لديه أقل عدد من الاتصالات النشطة
	server := lb.getLeastConnectionsServer()
	if server == nil {
		http.Error(w, "No healthy servers available", http.StatusServiceUnavailable)
		return
	}

	// زيادة عدد الاتصالات النشطة
	lb.mu.Lock()
	server.ActiveConnections++
	lb.mu.Unlock()

	// تسجيل عملية التوجيه
	lb.logger.Printf("Redirecting request to %s (Active connections: %d)\n", server.Name, server.ActiveConnections)

	// توجيه الطلب إلى الخادم
	targetURL, err := url.Parse(server.URL)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)

	// تقليل عدد الاتصالات بعد الانتهاء من معالجة الطلب
	lb.mu.Lock()
	server.ActiveConnections--
	lb.mu.Unlock()
}

// Function to check server health at intervals
func (lb *LoadBalancer) HealthCheck(interval time.Duration) {
	for {
		// قراءة القيمة من config
		for i := range lb.servers {
			resp, err := http.Get(lb.servers[i].URL + "/healthcheck")
			if err != nil || resp.StatusCode != http.StatusOK {
				lb.logger.Printf("Server %s is DOWN\n", lb.servers[i].Name)
				lb.servers[i].Healthy = false
			} else {
				lb.servers[i].Healthy = true
				lb.logger.Printf("Server %s is UP\n", lb.servers[i].Name)
			}
		}
		// الانتظار للمدة المحددة قبل التحقق مرة أخرى
		time.Sleep(interval)
	}
}

func main() {
	// Load configuration
	configData, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// تحويل HealthCheckInterval إلى time.Duration
	interval, err := time.ParseDuration(configData.HealthCheckInterval)
	if err != nil {
		log.Fatalf("Error parsing HealthCheckInterval: %v", err)
	}

	// Create the Load Balancer
	loadBalancer := NewLoadBalancer(configData.Servers)

	// Start health checks in a separate Goroutine
	go loadBalancer.HealthCheck(interval)

	// Start the Load Balancer server
	log.Printf("Load Balancer is running on port %s\n", configData.ListenPort)
	log.Fatal(http.ListenAndServe(configData.ListenPort, loadBalancer))
}
