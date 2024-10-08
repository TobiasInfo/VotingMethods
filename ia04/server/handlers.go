package server

import (
	//"encoding/json"
	//"ia04/agt"
	"ia04/comsoc"
	"net/http"
	"sync"
)

var (
	profile   comsoc.Profile
	voting    bool
	voteMutex sync.Mutex
)

func NewVoteHandler(w http.ResponseWriter, r *http.Request) {
	voteMutex.Lock()
	defer voteMutex.Unlock()

	profile = comsoc.Profile{}
	voting = true
	w.WriteHeader(http.StatusCreated)
}

// func VoteHandler(w http.ResponseWriter, r *http.Request) {
// 	voteMutex.Lock()
// 	defer voteMutex.Unlock()

// 	if !voting {
// 		http.Error(w, "Voting has ended", http.StatusForbidden)
// 		return
// 	}

// 	var agent agt.Agent
// 	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	profile = append(profile, agent.Prefs)
// 	w.WriteHeader(http.StatusOK)
// }

// func ResultHandler(w http.ResponseWriter, r *http.Request) {
// 	voteMutex.Lock()
// 	defer voteMutex.Unlock()

// 	if voting {
// 		http.Error(w, "Voting is still ongoing", http.StatusForbidden)
// 		return
// 	}

// 	bestAlts, err := comsoc.MajoritySCF(profile)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(bestAlts)
// }

func FinishHandler(w http.ResponseWriter, r *http.Request) {
	voteMutex.Lock()
	defer voteMutex.Unlock()

	voting = false
	w.WriteHeader(http.StatusOK)
}