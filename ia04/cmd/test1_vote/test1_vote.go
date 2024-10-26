package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	ballotData := map[string]interface{}{
		"rule":      "majority",
		"deadline":  "2020-10-09T23:05:08+02:00",
		"voter-ids": []string{"ag_id1", "ag_id2", "ag_id3"},
		"#alts":     4,
		"tie-break": []int{2, 4, 3, 1},
	}
	if !sendPostRequest("http://localhost:8080/new_ballot", ballotData) {
		return
	}

	votes := []map[string]interface{}{
		{
			"agent-id":  "ag_id1",
			"ballot-id": "scrutin1",
			"prefs":     []int{4, 2, 3, 1},
		},
		{
			"agent-id":  "ag_id2",
			"ballot-id": "scrutin1",
			"prefs":     []int{1, 2, 3, 4},
		},
		{
			"agent-id":  "ag_id3",
			"ballot-id": "scrutin1",
			"prefs":     []int{4, 2, 3, 1},
		},
	}

	for _, vote := range votes {
		if !sendPostRequest("http://localhost:8080/vote", vote) {
			return
		}
	}

	resultData := map[string]interface{}{
		"ballot-id": "scrutin1",
	}
	result := sendPostRequestWithResponse("http://localhost:8080/result", resultData)
	if result != nil {
		fmt.Println("Final Result:", string(result))
	} else {
		fmt.Println("Failed to retrieve result.")
	}
}

func sendPostRequest(url string, data map[string]interface{}) bool {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return false
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request to %s failed: %v\n", url, err)
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

func sendPostRequestWithResponse(url string, data map[string]interface{}) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return nil
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request to %s failed: %v\n", url, err)
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
