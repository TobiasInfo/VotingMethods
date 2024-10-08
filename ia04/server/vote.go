package server

import (
    "encoding/json"
    "net/http"
    "time"
)

// Vote structure to store the vote information
type Vote struct {
    AgentID  string   `json:"agent-id"`
    BallotID string   `json:"ballot-id"`
    Prefs    []int    `json:"prefs"`
    Options  []int    `json:"options,omitempty"`  // Optional field
}

var votes = map[string]map[string]Vote{}  // Mapping of BallotID -> AgentID -> Vote

func voteHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var vote Vote
    err := json.NewDecoder(r.Body).Decode(&vote)
    if err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }

    if vote.AgentID == "" || vote.BallotID == "" || len(vote.Prefs) == 0 {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // Check if the ballot exists
    ballot, exists := ballots[vote.BallotID]
    if !exists {
        http.Error(w, "Ballot not found", http.StatusBadRequest)
        return
    }

    // Check if the voting deadline has passed
    if time.Now().After(ballot.Deadline) {
        http.Error(w, "Voting deadline has passed", http.StatusServiceUnavailable)  // 503
        return
    }

    // Check if the agent has already voted
    if _, voted := votes[vote.BallotID][vote.AgentID]; voted {
        http.Error(w, "Agent has already voted", http.StatusForbidden)  // 403
        return
    }

	// if len(vote.Options) > 0 {
	// 	for _, option := range vote.Options {
	// 		if !checkProfileAlternative(option, ballot.NumAlts) {
	// 			http.Error(w, "Invalid option value", http.StatusBadRequest)
	// 			return
	// 		}
	// 	}
	// }

    // Record the vote
    if votes[vote.BallotID] == nil {
        votes[vote.BallotID] = map[string]Vote{}
    }
    votes[vote.BallotID][vote.AgentID] = vote

    // Return success
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Vote successfully recorded"))
}
