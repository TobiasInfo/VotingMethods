package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

var nbOfRequests = 0

type Ballot struct {
	Rule     string    `json:"rule"`
	Deadline time.Time `json:"deadline"`
	VoterIDs []string  `json:"voter-ids"`
	NumAlts  int       `json:"#alts"`
	TieBreak []int     `json:"tie-break"`
	BallotID string    `json:"ballot-id"`
}

var ballots = map[string]Ballot{}

// TODO : Add vote that are implemented
var supportedRules = map[string]bool{
	"majority": true,
	"borda":    true,
	"approval": true,
	"stv":      true,
}

func newBallotHandler(w http.ResponseWriter, r *http.Request) {
	// TODO : Check if we need mutex or channel to increment nbOfRequests
	nbOfRequests++
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

	if ballot.Rule == "" || ballot.Deadline.IsZero() || len(ballot.VoterIDs) == 0 || ballot.NumAlts == 0 || len(ballot.TieBreak) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if !supportedRules[ballot.Rule] {
		http.Error(w, "Voting rule not implemented", http.StatusNotImplemented)
		return
	}

	if ballot.NumAlts < 2 {
		http.Error(w, "Invalid number of alternatives", http.StatusBadRequest)
		return
	}

	if len(ballot.TieBreak) != ballot.NumAlts {
		http.Error(w, "Invalid tie-breaking function", http.StatusBadRequest)
		return
	}

	ballot.BallotID = "scrutin" + strconv.Itoa(nbOfRequests)

	ballots[ballot.BallotID] = ballot

	response := map[string]string{
		"ballot-id": ballot.BallotID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
