package main

import (
	"github.com/TobiasInfo/SystemeMultiAgents/server"
	"log"
	"net/http"
)

func main() {
	router := server.NewRouter()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
