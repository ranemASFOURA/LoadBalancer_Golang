package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func sendRequest(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		log.Printf("Request %d failed: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Request %d -> Response: %s\n", id, body)
}

func main() {
	var wg sync.WaitGroup
	numRequests := 5
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go sendRequest(i, &wg)
	}

	wg.Wait()
	fmt.Println("All requests completed.")
}
