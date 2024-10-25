package server

import (
	"encoding/json"
	"fmt"
	"ia04/comsoc"

	// "log"
	"net/http"
	// "time"
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
	thresholds := make([]int, 0, len(ballotVotes))
	if ballot.Rule == "approval" {
		for _, vote := range ballotVotes {
			if len(vote.Options) == 1 {
				thresholds = append(thresholds, vote.Options[0])
			} else {
				return 0, nil, fmt.Errorf("approval vote must have only one option")
			}
		}
	}
	for _, vote := range ballotVotes {
		// Convert vote.Prefs to comsoc.Alternative
		comsocVote := make([]comsoc.Alternative, len(vote.Prefs))
		for i, alt := range vote.Prefs {
			comsocVote[i] = comsoc.Alternative(alt)
		}
		profile = append(profile, comsocVote)
	}

	var winner comsoc.Alternative
	var ranking []comsoc.Alternative
	var err error

	// var swf func(p comsoc.Profile) (comsoc.Count, error)
	var tieBreaker func([]comsoc.Alternative) (comsoc.Alternative, error)
	var scf func(comsoc.Profile) ([]comsoc.Alternative, error)
	var scfapproval func(comsoc.Profile, []int) ([]comsoc.Alternative, error)
	var scffactory func(comsoc.Profile) (comsoc.Alternative, error)

	// Convert TieBreak from int to comsoc.Alternative
	tieBreakSlice := make([]comsoc.Alternative, len(ballot.TieBreak))
	for i, alt := range ballot.TieBreak {
		tieBreakSlice[i] = comsoc.Alternative(alt)
	}
	tieBreaker = comsoc.TieBreakFactory(tieBreakSlice)

	switch ballot.Rule {
	case "majority":
		scf = comsoc.MajoritySCF
	case "borda":
		scf = comsoc.BordaSCF
	case "approval":
		scfapproval = comsoc.ApprovalSCF
	case "condorcet":
		scf = comsoc.CondorcetWinner

	case "copeland":
		scf = comsoc.CopelandSCF

	default:
		return 0, ranking, fmt.Errorf("unsupported voting rule")
	}
	if ballot.Rule == "approval" {
		alts, err := scfapproval(profile, thresholds)
		if err != nil {
			return 0, ranking, fmt.Errorf("failed to get alternatives from SCF: %w", err)
		}
		if len(alts) == 1 {
			return alts[0], ranking, nil
		}
		alt, err := tieBreaker(alts)
		if err != nil {
			return 0, ranking, fmt.Errorf("failed to get best alternative from tie-breaking function: %w", err)
		}
		return alt, ranking, nil
	}
	scffactory = comsoc.SCFFactory(scf, tieBreaker)
	winner, err = scffactory(profile)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to compute the winner: %w", err)
	}
	ranking, err = []comsoc.Alternative{}, nil

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
	// if time.Now().Before(ballot.Deadline) {
	// 	http.Error(w, "Voting still ongoing", http.StatusTooEarly)
	// 	return
	// }

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
