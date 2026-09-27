package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/noah-philip/skyos/internal"
)

// Creates an HTTP server
func main() {

	//Command line flag
	NodeID := flag.String(
		"id",
		"drone-1",
		"unique ID for this drone",
	)

	//Port
	port := flag.String(
		"port",
		"8081",
		"HTTP port for this drone",
	)

	//Reads in arguments passed in command line & fills in NodeID and port
	flag.Parse()

	//Creates Drone 1's local peer registry
	registry := internal.NewPeerRegistry()

	mux := http.NewServeMux()

	//HTTP request router, when a specific HTTP request arrives at this path, what function should handle this?
	mux.HandleFunc(
		"/v1/heartbeat",
		registry.HandleHeartbeat,
	)

	address := ":" + *port

	//Prints message
	log.Printf(
		"droned %s listening on http://localhost%s",
		*NodeID,
		address,
	)

	//Actually starts the server. Opens specified port, waits for HTTP request, sends requests to mux, keeps running till error or shutdown
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal(err)
	}
}
