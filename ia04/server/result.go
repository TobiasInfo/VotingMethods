package server

import (
	"encoding/json"
	"fmt"
	"ia04/comsoc"
	"net/http"
	"time"
)

// TODO : Implement Tie-Breaking, for now we will just return the first alternative
func computeResult(rule string, ballotVotes map[string]Vote, numAlts int) (comsoc.Alternative, []comsoc.Alternative, error) {
    switch rule {
    case "majority":
		bestAlts, err := comsoc.MajoritySCF(profile)
		if err != nil {
			return 0, nil, err
		} else {
            return bestAlts[0], bestAlts, nil
		}
    case "borda":
		bestAlts, err := comsoc.BordaSCF(profile)
		if err != nil {
			return 0, nil, err
		} else {
			return bestAlts[0], bestAlts, nil
		}
		
	// TODO : check how to use thresholds
    // case "approval":
	// 	bestAlts, err := comsoc.ApprovalSCF(profile)
	// 	if err != nil {
	// 		return 0, nil, err
	// 	} else {
	// 		return bestAlts[0], bestAlts, nil
	// 	}
	// case "stv":
	// 	bestAlts, err := comsoc.STVSCF(profile)
	// 	if err != nil {
	// 		return 0, nil, err
	// 	} else {
	// 		return bestAlts[0], bestAlts, nil
	// 	}
	case "copeland":
		bestAlts, err := comsoc.CopelandSCF(profile)
		if err != nil {
			return 0, nil, err
		} else {
			return bestAlts[0], bestAlts, nil
		}
	case "condorcet":
		bestAlts, err := comsoc.CondorcetWinner(profile)
		if err != nil {
			return 0, nil, err
		} else {
			return bestAlts[0], bestAlts, nil
		}
    default:
        return 0, nil, fmt.Errorf("unsupported voting rule: %s", rule)
    }
}


func resultHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        BallotID string `json:"ballot-id"`
    }
    
    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil || request.BallotID == "" {
        http.Error(w, "Invalid request body or missing ballot-id", http.StatusBadRequest)
        return
    }

    // Check if the ballot exists
    ballot, exists := ballots[request.BallotID]
    if !exists {
        http.Error(w, "Ballot not found", http.StatusNotFound)
        return
    }

    // Check if the voting deadline has passed
    if time.Now().Before(ballot.Deadline) {
        http.Error(w, "Voting still ongoing", http.StatusTooEarly)
        return
    }

    // Get the votes for the ballot
    ballotVotes, voted := votes[request.BallotID]
    if !voted || len(ballotVotes) == 0 {
        http.Error(w, "No votes have been cast for this ballot", http.StatusNotFound)
        return
    }

    // Compute the result based on the rule
    winner, ranking, err := computeResult(ballot.Rule, ballotVotes, ballot.NumAlts)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Build the response
    response := map[string]interface{}{
        "winner":  winner,
    }

    // If ranking exists, add it to the response
    if len(ranking) > 0 {
        response["ranking"] = ranking
    }

    // Return the result
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
