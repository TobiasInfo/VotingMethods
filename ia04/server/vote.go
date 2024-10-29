package server

import (
	"encoding/json"
	"net/http"
	"time"
)

// Vote structure to store the vote information
type Vote struct {
	AgentID  string `json:"agent-id"`
	BallotID string `json:"ballot-id"`
	Prefs    []int  `json:"prefs"`
	Options  []int  `json:"options,omitempty"` // Optional field
}

var votes = map[string]map[string]Vote{} // Mapping of BallotID -> AgentID -> Vote

// Channel for processing votes
var voteChannel = make(chan VoteRequest)

type VoteRequest struct {
	Vote     Vote
	RespChan chan Response
}

type Response struct {
	StatusCode int
	Message    string
}

// Initialize the vote processor
func init() {
	go processVotes()
}

// Function to process votes from the channel
func processVotes() {
	for req := range voteChannel {
		vote := req.Vote
		response := Response{}

		// Check if the ballot exists
		ballot, exists := ballots[vote.BallotID]
		if !exists {
			response = Response{StatusCode: http.StatusBadRequest, Message: "Ballot not found"}
			req.RespChan <- response
			continue
		}

		// Check if the agent is a valid voter
		validVoter := false
		for _, voterID := range ballot.VoterIDs {
			if vote.AgentID == voterID {
				validVoter = true
				break
			}
		}
		if !validVoter {
			response = Response{StatusCode: http.StatusBadRequest, Message: "Agent is not a valid voter"}
			req.RespChan <- response
			continue
		}

		// Check if the voting deadline has passed
		if time.Now().After(ballot.Deadline) {
			response = Response{StatusCode: http.StatusServiceUnavailable, Message: "Voting deadline has passed"}
			req.RespChan <- response
			continue
		}

		// Check if the agent has already voted
		if _, voted := votes[vote.BallotID][vote.AgentID]; voted {
			response = Response{StatusCode: http.StatusForbidden, Message: "Agent has already voted"}
			req.RespChan <- response
			continue
		}

		// Record the vote
		if votes[vote.BallotID] == nil {
			votes[vote.BallotID] = map[string]Vote{}
		}
		votes[vote.BallotID][vote.AgentID] = vote

		// Send the success response
		response = Response{StatusCode: http.StatusOK, Message: "Vote successfully recorded"}
		req.RespChan <- response
	}
}

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

	// Channel to receive the response
	responseChannel := make(chan Response)

	// Send the vote request to the channel
	voteChannel <- VoteRequest{Vote: vote, RespChan: responseChannel}

	// Wait for the response
	response := <-responseChannel

	// Return the response
	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Message))
}
