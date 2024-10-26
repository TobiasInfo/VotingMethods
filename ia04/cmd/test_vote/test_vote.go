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

func new_ballot_request(url string, data map[string]interface{}) bool {
	resp, err := sendPostRequest(url, data)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("Request to %s failed with status: %s\n", url, resp.Status)
		return false
	}

	fmt.Printf("Request to %s succeeded with status: %s\n", url, resp.Status)
	return true
}

func vote_request(ballotID string) {
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
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			fmt.Printf("Request to %s failed with status: %s\n", "http://localhost:8080/vote", resp.Status)
			return
		}

		fmt.Printf("Request to %s succeeded with status: %s\n", "http://localhost:8080/vote", resp.Status)
	}
}

func result_request(url string, data map[string]interface{}) []byte {
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
	var wg_create_ballot sync.WaitGroup
	wg_create_ballot.Add(100)
	for i := 1; i <= 100; i++ {
		go func(i int) {
			defer wg_create_ballot.Done()
			ballotData := map[string]interface{}{
				"rule":      "majority",
				"deadline":  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"voter-ids": []string{"ag_id1", "ag_id2", "ag_id3"},
				"#alts":     4,
				"tie-break": []int{4, 1, 3, 2},
			}
			if !new_ballot_request("http://localhost:8080/new_ballot", ballotData) {
				return
			}
		}(i)
	}

	wg_create_ballot.Wait()
	fmt.Println("All new ballot requests completed")

	var wg_vote sync.WaitGroup
	wg_vote.Add(100)

	for i := 1; i <= 100; i++ {
		go func(i int) {
			defer wg_vote.Done()
			ballotID := fmt.Sprintf("scrutin%d", i)
			vote_request(ballotID)
		}(i)
	}

	wg_vote.Wait()
	fmt.Println("All vote requests completed")

	var wg_get_results sync.WaitGroup
	wg_get_results.Add(100)

	for i := 1; i <= 100; i++ {
		go func(i int) {
			defer wg_get_results.Done()
			ballotID := fmt.Sprintf("scrutin%d", i)
			resultData := map[string]interface{}{
				"ballot-id": ballotID,
			}
			result := result_request("http://localhost:8080/result", resultData)
			if result != nil {
				fmt.Printf("Final Result for %s: %s\n", ballotID, string(result))
			} else {
				fmt.Printf("Failed to retrieve result for %s.\n", ballotID)
			}
		}(i)
	}

	wg_get_results.Wait()
	fmt.Println("All result retrievals completed")
}
