package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	const numRequests = 100
	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func(i int) {
			defer wg.Done()
			sendRequest(i)
		}(i)
	}

	// Wait for all requests to finish
	wg.Wait()
	fmt.Println("All requests completed")
}

func sendRequest(id int) {
	// Sample ballot JSON data
	ballotData := []byte(fmt.Sprintf(`{
		"rule": "majority",
		"deadline": "%s",
		"voter-ids": ["ag_id1", "ag_id2", "ag_id3"],
		"#alts": 3,
		"tie-break": [1, 2, 3]
	}`, time.Now().Add(24*time.Hour).Format(time.RFC3339)))

	// Create the POST request
	resp, err := http.Post("http://localhost:8080/new_ballot", "application/json", bytes.NewBuffer(ballotData))
	if err != nil {
		fmt.Printf("Request %d failed: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("Request %d succeeded\n", id)
	} else {
		fmt.Printf("Request %d failed with status: %s\n", id, resp.Status)
	}
}
