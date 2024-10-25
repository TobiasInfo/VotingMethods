package server

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/newvote", NewVoteHandler).Methods("POST")
	//router.HandleFunc("/vote", VoteHandler).Methods("POST")
	//router.HandleFunc("/result", ResultHandler).Methods("GET")
	router.HandleFunc("/finish", FinishHandler).Methods("POST")

	router.HandleFunc("/new_ballot", createBallotHandler).Methods("POST")
	router.HandleFunc("/vote", voteHandler).Methods("POST")
	router.HandleFunc("/result", resultHandler).Methods("POST")
	return router
}
