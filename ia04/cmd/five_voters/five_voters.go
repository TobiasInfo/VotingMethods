package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ia04/agt"
)

func main() {
	// Start a new voting session
	_, err := http.Post("http://localhost:8080/newvote", "application/json", nil)
	if err != nil {
		log.Fatalf("Failed to start new vote: %v", err)
	}

	// Create some agents
	agent1 := &agt.Agent{ID: 1, Name: "Alice", Prefs: []agt.Alternative{1, 2, 3}}
	agent2 := &agt.Agent{ID: 2, Name: "Bob", Prefs: []agt.Alternative{2, 3, 1}}
	agent3 := &agt.Agent{ID: 3, Name: "Charlie", Prefs: []agt.Alternative{1, 3, 2}}
	agent4 := &agt.Agent{ID: 4, Name: "Dave", Prefs: []agt.Alternative{3, 1, 2}}
	agent5 := &agt.Agent{ID: 5, Name: "Eve", Prefs: []agt.Alternative{2, 1, 3}}

	// Submit votes
	agents := []*agt.Agent{agent1, agent2, agent3, agent4, agent5}
	for _, agent := range agents {
		voteData, _ := json.Marshal(agent)
		_, err := http.Post("http://localhost:8080/vote", "application/json", bytes.NewBuffer(voteData))
		if err != nil {
			log.Fatalf("Failed to submit vote: %v", err)
		}
	}

	// Finish voting
	_, err = http.Post("http://localhost:8080/finish", "application/json", nil)
	if err != nil {
		log.Fatalf("Failed to finish voting: %v", err)
	}

	// Wait a moment for the server to process
	time.Sleep(1 * time.Second)

	// Get the result
	resp, err := http.Get("http://localhost:8080/result")
	if err != nil {
		log.Fatalf("Failed to get result: %v", err)
	}
	defer resp.Body.Close()

	var result []agt.Alternative
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode result: %v", err)
	}

	log.Printf("Best alternatives according to Majority SCF: %v", result)
}
