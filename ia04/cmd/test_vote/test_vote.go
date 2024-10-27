package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func sendPostRequest(url string, data map[string]interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %v", url, err)
	}

	return resp, nil
}

func newBallotRequest(url string, data map[string]interface{}, done chan<- bool) {
	resp, err := sendPostRequest(url, data)
	if err != nil {
		fmt.Println(err)
		done <- false
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("Request to %s failed with status: %s\n", url, resp.Status)
		done <- false
		return
	}

	fmt.Printf("Request to %s succeeded with status: %s\n", url, resp.Status)
	done <- true
}

func voteRequest(ballotID string, done chan<- bool) {
	votes := []map[string]interface{}{
		{
			"agent-id":  "ag_id1",
			"ballot-id": ballotID,
			"prefs":     []int{4, 2, 3, 1},
			"options":   []int{1},
		},
		{
			"agent-id":  "ag_id2",
			"ballot-id": ballotID,
			"prefs":     []int{1, 4, 3, 2},
			"options":   []int{2},
		},
		{
			"agent-id":  "ag_id3",
			"ballot-id": ballotID,
			"prefs":     []int{1, 2, 3, 4},
			"options":   []int{1},
		},
	}

	for _, vote := range votes {
		resp, err := sendPostRequest("http://localhost:8080/vote", vote)
		if err != nil {
			fmt.Println(err)
			done <- false
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			fmt.Printf("Request to %s failed with status: %s\n", "http://localhost:8080/vote", resp.Status)
			done <- false
			return
		}

		fmt.Printf("Request to %s succeeded with status: %s\n", "http://localhost:8080/vote", resp.Status)
	}
	done <- true
}

func resultRequest(url string, data map[string]interface{}) []byte {
	resp, err := sendPostRequest(url, data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request to %s failed with status: %s\n", url, resp.Status)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v\n", err)
		return nil
	}
	return body
}

func main() {
	// Channel for ballot creation
	createBallotDone := make(chan bool)
	var wgCreateBallot sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wgCreateBallot.Add(1)
		go func(i int) {
			defer wgCreateBallot.Done()
			ballotData := map[string]interface{}{
				"rule":      "approval",
				"deadline":  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"voter-ids": []string{"ag_id1", "ag_id2", "ag_id3"},
				"#alts":     4,
				"tie-break": []int{4, 1, 3, 2},
			}
			newBallotRequest("http://localhost:8080/new_ballot", ballotData, createBallotDone)
		}(i)
	}

	// Wait for all ballots to be created
	go func() {
		wgCreateBallot.Wait()
		close(createBallotDone)
	}()

	// Wait for ballot creation to complete
	for success := range createBallotDone {
		if !success {
			fmt.Println("Ballot creation failed for one or more ballots.")
		}
	}
	fmt.Println("All new ballot requests completed")

	// Channel for vote requests
	voteDone := make(chan bool)
	var wgVote sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wgVote.Add(1)
		go func(i int) {
			defer wgVote.Done()
			ballotID := fmt.Sprintf("scrutin%d", i)
			voteRequest(ballotID, voteDone)
		}(i)
	}

	// Wait for all votes to be cast
	go func() {
		wgVote.Wait()
		close(voteDone)
	}()

	// Wait for vote requests to complete
	for success := range voteDone {
		if !success {
			fmt.Println("Voting failed for one or more ballots.")
		}
	}
	fmt.Println("All vote requests completed")

	// Channel for result requests
	resultDone := make(chan bool)
	var wgResult sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wgResult.Add(1)
		go func(i int) {
			defer wgResult.Done()
			ballotID := fmt.Sprintf("scrutin%d", i)
			resultData := map[string]interface{}{
				"ballot-id": ballotID,
			}
			result := resultRequest("http://localhost:8080/result", resultData)
			if result != nil {
				fmt.Printf("Final Result for %s: %s\n", ballotID, string(result))
			} else {
				fmt.Printf("Failed to retrieve result for %s.\n", ballotID)
			}
			resultDone <- result != nil
		}(i)
	}

	// Wait for all results to be retrieved
	go func() {
		wgResult.Wait()
		close(resultDone)
	}()

	// Collect result request outcomes
	for success := range resultDone {
		if !success {
			fmt.Println("Result retrieval failed for one or more ballots.")
		}
	}
	fmt.Println("All result retrievals completed")
}
