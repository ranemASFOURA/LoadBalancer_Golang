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

var requestCount int
var requestMu sync.Mutex

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run server.go <ServerName>")
	}

	serverName := os.Args[1]

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

	if !found {
		log.Fatalf("Server %s not found in config.json", serverName)
	}

	// استخراج المنفذ فقط من الـ URL
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		log.Fatalf("Invalid URL for server %s: %v", serverName, err)
	}

	port := parsedURL.Port()
	if port == "" {
		log.Fatalf("No port found in URL for server %s", serverName)
	}

	// التحقق من طلب Health Check أولاً
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		// التحقق من طلب Health Check
		fmt.Printf("%s received a HealthCheck request.\n", serverName)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HealthCheck OK"))
	})

	// التعامل مع باقي الطلبات
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestMu.Lock()
		requestCount++
		fmt.Printf("%s received a request. Total requests handled: %d\n", serverName, requestCount)
		requestMu.Unlock()
		fmt.Fprintf(w, "Hello from %s\n", serverName)
	})

	log.Printf("%s is running on port %s\n", serverName, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
