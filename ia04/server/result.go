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

	var winner comsoc.Alternative = 0
	ranking := make([]comsoc.Alternative, 0)
	var err error
	// var swf func(p comsoc.Profile) (comsoc.Count, error)
	var tieBreaker func([]comsoc.Alternative) (comsoc.Alternative, error)
	var scf func(comsoc.Profile) ([]comsoc.Alternative, error)
	var swf func(comsoc.Profile) (comsoc.Count, error)
	var scfapproval func(comsoc.Profile, []int) ([]comsoc.Alternative, error)
	var swfapproval func(comsoc.Profile, []int) (comsoc.Count, error)
	var scffactory func(comsoc.Profile) (comsoc.Alternative, error)
	var swffactory func(comsoc.Profile) ([]comsoc.Alternative, error)
	// Convert TieBreak from int to comsoc.Alternative
	tieBreakSlice := make([]comsoc.Alternative, len(ballot.TieBreak))
	for i, alt := range ballot.TieBreak {
		tieBreakSlice[i] = comsoc.Alternative(alt)
	}
	tieBreaker = comsoc.TieBreakFactory(tieBreakSlice)

	switch ballot.Rule {
	case "majority":
		scf = comsoc.MajoritySCF
		swf = comsoc.MajoritySWF
	case "borda":
		scf = comsoc.BordaSCF
		swf = comsoc.BordaSWF
	case "approval":
		scfapproval = comsoc.ApprovalSCF
		swfapproval = comsoc.ApprovalSWF
	case "condorcet":
		scf = comsoc.CondorcetWinner
		swf = nil
	case "copeland":
		scf = comsoc.CopelandSCF
		swf = comsoc.CopelandSWF

	default:
		return 0, ranking, fmt.Errorf("unsupported voting rule")
	}
	if ballot.Rule == "approval" {
		alts, err := scfapproval(profile, thresholds)
		if err != nil {
			return winner, ranking, fmt.Errorf("failed to get alternatives from SCF: %w", err)
		}
		if len(alts) == 1 {
			winner = alts[0]
			return winner, ranking, nil
		}
		alt, err := tieBreaker(alts)
		if err != nil {
			return winner, ranking, fmt.Errorf("failed to get best alternative from tie-breaking function: %w", err)
		}
		winner = alt
		count, err := swfapproval(profile, thresholds)
		// log.Println("count", count)
		if err != nil {
			return winner, ranking, fmt.Errorf("failed to get count from SWF: %w", err)
		}
		for {
			bestAlt := comsoc.MaxCount(count)
			if len(bestAlt) == 0 {
				break
			}
			if len(bestAlt) == 1 {
				ranking = append(ranking, bestAlt[0])
				delete(count, bestAlt[0])
			} else {
				alt, err := tieBreaker(bestAlt)
				if err != nil {
					return winner, ranking, fmt.Errorf("failed to get best alternative from tie-breaking function: %w", err)
				}
				ranking = append(ranking, alt)
				delete(count, alt)
			}
		}
		// log.Println("ranking", ranking)
		return winner, ranking, nil
	}
	scffactory = comsoc.SCFFactory(scf, tieBreaker)
	winner, err = scffactory(profile)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to compute the winner: %w", err)
	}
	if swf == nil {
		ranking = append(ranking, winner)
		return winner, ranking, nil
	}
	// count, _ := swf(profile)
	// log.Println("count", count)
	swffactory = comsoc.SWFFactory(swf, tieBreaker)
	ranking, err = swffactory(profile)
	// log.Println("ranking", ranking)
	if err != nil {
		return winner, ranking, fmt.Errorf("failed to compute ranking: %w", err)
	}
	return winner, ranking, nil
}

var resultChannel = make(chan ResultRequest)

type ResultRequest struct {
	BallotID string
	RespChan chan ResultResponse
}

type ResultResponse struct {
	StatusCode int
	Message    string
	Data       map[string]interface{}
}

func init() {
	go processResults()
}

// Function to process result requests from the channel
func processResults() {
	for req := range resultChannel {
		response := ResultResponse{}

		// Check if the ballot exists
		ballot, exists := ballots[req.BallotID]
		if !exists {
			response = ResultResponse{StatusCode: http.StatusNotFound, Message: "Ballot not found"}
			req.RespChan <- response
			continue
		}

		// Check if voting is complete (optional deadline check commented out)
		// if time.Now().Before(ballot.Deadline) {
		// 	response = ResultResponse{StatusCode: http.StatusTooEarly, Message: "Voting still ongoing"}
		// 	req.RespChan <- response
		// 	continue
		// }

		// Retrieve the votes for the ballot
		ballotVotes, voted := votes[req.BallotID]
		if !voted || len(ballotVotes) == 0 {
			response = ResultResponse{StatusCode: http.StatusNotFound, Message: "No votes have been cast for this ballot"}
			req.RespChan <- response
			continue
		}

		// Compute the result based on the rule
		winner, ranking, err := computeResult(ballot, ballotVotes)
		if err != nil {
			response = ResultResponse{StatusCode: http.StatusInternalServerError, Message: err.Error()}
			req.RespChan <- response
			continue
		}

		// Prepare the response data
		data := map[string]interface{}{
			"winner": winner,
		}
		if len(ranking) > 0 {
			data["ranking"] = ranking
		}

		// Send the success response with data
		response = ResultResponse{StatusCode: http.StatusOK, Message: "Result computed successfully", Data: data}
		req.RespChan <- response
	}
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

	// Channel to receive the response
	responseChannel := make(chan ResultResponse)

	// Send the result request to the channel
	resultChannel <- ResultRequest{BallotID: request.BallotID, RespChan: responseChannel}

	// Wait for the response
	response := <-responseChannel

	// Return the response
	w.WriteHeader(response.StatusCode)
	if response.StatusCode == http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.Data)
	} else {
		w.Write([]byte(response.Message))
	}
}
