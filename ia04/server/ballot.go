package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var nbOfRequests = 0
var requestMutex sync.Mutex // Mutex to protect nbOfRequests

type Ballot struct {
	Rule     string    `json:"rule"`
	Deadline time.Time `json:"deadline"`
	VoterIDs []string  `json:"voter-ids"`
	NumAlts  int       `json:"#alts"`
	TieBreak []int     `json:"tie-break"`
	BallotID string    `json:"ballot-id"`
}

var ballots = map[string]Ballot{}
var ballotRequests = make(chan BallotRequest) // Channel for ballot creation requests

// Supported voting rules
var supportedRules = map[string]bool{
	"majority": true,
	"borda":    true,
	"approval": true,
	"stv":      true,
	"kemeny":   true, // Added "kemeny" as supported
}

type BallotRequest struct {
	Ballot Ballot
	Resp   chan string // Channel to send the response back to the handler
}

func init() {
	// Start a goroutine to process ballot requests
	go func() {
		for request := range ballotRequests {
			createBallot(request)
		}
	}()
}

func createBallotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var ballot Ballot
	err := json.NewDecoder(r.Body).Decode(&ballot)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Validate the fields
	if !supportedRules[ballot.Rule] {
		http.Error(w, "Voting rule not implemented", http.StatusNotImplemented)
		return
	}

	if len(ballot.VoterIDs) == 0 {
		http.Error(w, "Voter IDs are required", http.StatusBadRequest)
		return
	}

	if ballot.NumAlts < 2 {
		http.Error(w, "Invalid number of alternatives", http.StatusBadRequest)
		return
	}

	if len(ballot.TieBreak) != ballot.NumAlts {
		http.Error(w, "Tie-break length must match number of alternatives", http.StatusBadRequest)
		return
	}

	// Create a response channel for the ballot creation request
	responseChan := make(chan string)
	ballotRequests <- BallotRequest{Ballot: ballot, Resp: responseChan}

	// Wait for the response from the ballot creation goroutine
	ballotID := <-responseChan

	response := map[string]string{
		"ballot-id": ballotID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func createBallot(request BallotRequest) {
	// Increment the nbOfRequests safely using the mutex
	requestMutex.Lock()
	nbOfRequests++
	request.Ballot.BallotID = "scrutin" + strconv.Itoa(nbOfRequests)
	requestMutex.Unlock()

	// Store the ballot in the map
	ballots[request.Ballot.BallotID] = request.Ballot

	// Send the response back to the handler
	request.Resp <- request.Ballot.BallotID
}
