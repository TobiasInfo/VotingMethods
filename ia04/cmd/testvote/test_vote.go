package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	const numRequests = 100
	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Define the valid agent IDs and ballots
	agentIDs := []string{"ag_id1", "ag_id2", "ag_id3", "ag_id4"}

	for i := 0; i < numRequests; i++ {
		go func(i int) {
			defer wg.Done()
			sendVoteRequest(agentIDs[i%len(agentIDs)], i%99+1) // Use agent IDs and rotate through ballots 1 to 99
		}(i)
	}

	// Wait for all requests to finish
	wg.Wait()
	fmt.Println("All requests completed")
}

func sendVoteRequest(agentID string, ballotNumber int) {
	// Sample vote JSON data
	voteData := []byte(fmt.Sprintf(`{
		"agent-id": "%s",
		"ballot-id": "scrutin%d",
		"prefs": [1, 2, 3, 4],
		"options": [3]
	}`, agentID, ballotNumber)) // Use the specified agent ID and ballot number

	// Create the POST request
	resp, err := http.Post("http://localhost:8080/vote", "application/json", bytes.NewBuffer(voteData))
	if err != nil {
		fmt.Printf("Vote request for agent %s on ballot scrutin%d failed: %v\n", agentID, ballotNumber, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Vote request for agent %s on ballot scrutin%d succeeded\n", agentID, ballotNumber)
	} else {
		fmt.Printf("Vote request for agent %s on ballot scrutin%d failed with status: %s\n", agentID, ballotNumber, resp.Status)
	}
}
