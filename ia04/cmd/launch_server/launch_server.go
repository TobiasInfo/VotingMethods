package main

import (
	"log"
	"net/http"

	"github.com/TobiasInfo/SystemeMultiAgents/ia04/server"
)

func main() {
	router := server.NewRouter()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
