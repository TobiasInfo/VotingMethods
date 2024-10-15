package server

import (
	"encoding/json"
	"fmt"
	"ia04/comsoc"
	"net/http"
	"time"
)



// Compute the result based on the voting rule and the votes received
func computeResult(ballot Ballot, ballotVotes map[string]Vote) (comsoc.Alternative, []comsoc.Alternative, error) {
    // Recover the vote method from the ballot
    _, isRuleSupported := supportedRules[ballot.Rule]
    if !isRuleSupported {
		return 0, nil, fmt.Errorf("unsupported voting rule")
    }

    // Create the profile from the votes
    profile := comsoc.Profile{}
    for _, vote := range ballotVotes {
        // Convert vote.Prefs to comsoc.Alternative
        comsocVote := make([]comsoc.Alternative, len(vote.Prefs))
        for i, alt := range vote.Prefs {
            comsocVote[i] = comsoc.Alternative(alt)
        }
    }

    var winner comsoc.Alternative
    var ranking []comsoc.Alternative
    var err error

    var swf func(p comsoc.Profile) (comsoc.Count, error)
	var tieBreaker func([]comsoc.Alternative) (comsoc.Alternative, error)
    
    // Convert TieBreak from int to comsoc.Alternative
    tieBreakSlice := make([]comsoc.Alternative, len(ballot.TieBreak))
    for i, alt := range ballot.TieBreak {
        tieBreakSlice[i] = comsoc.Alternative(alt)
    }

    tieBreaker = comsoc.TieBreakFactory(tieBreakSlice)

    switch ballot.Rule {
        case "majority":
            swf = comsoc.MajoritySWF
        case "borda":
            swf = comsoc.BordaSWF
        // case "approval":
        //     best, err := comsoc.ApprovalSWF(profile, ballot.TieBreak)
        //     if err != nil {
        //         return 0, nil, fmt.Errorf("failed to compute the winner: %w", err)
        //     }
        //     tieBreaker = comsoc.TieBreakFactory(ballot.TieBreak)
        //     return 0, best, nil
        
        // case "condorcet":
        //     swf = comsoc.CondorcetWinner

        case "copeland":
            swf = comsoc.CopelandSWF

        default:
            return 0, nil, fmt.Errorf("unsupported voting rule")
        }

        scf := comsoc.SWFFactory(swf, tieBreaker)
        ranking, err = scf(profile)
        if err != nil {
            return 0, nil, fmt.Errorf("failed to compute the winner: %w", err)
        }

        if len(ranking) > 0 {
            winner, err = tieBreaker(ranking)
        }
        return winner, ranking, nil
}

// Result handler to compute the result of a voting session
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
	winner, ranking, err := computeResult(ballot, ballotVotes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build the response
	response := map[string]interface{}{
		"winner": winner,
	}

	if len(ranking) > 0 {
		response["ranking"] = ranking
	}

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}