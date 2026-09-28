package main

import (
	"flag"
	"log"
	"net/http"
	"time"

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

	//Command line option for peerURL
	peerURL := flag.String(
		"peer",
		"", //Default value if you dont provide it
		"base URL of the peer drone",
	)
	battery := flag.Float64(
		"battery",
		100,
		"simulated battery percentage",
	)

	//Reads in arguments passed in command line & fills in NodeID and port
	flag.Parse()

	node := internal.Node{
		ID: *NodeID,
		Resources: internal.Resources{
			Battery: *battery,
		},
	}

	//Creates Drone 1's local peer registry
	registry := internal.NewPeerRegistry()

	//Creates Drone 1's http client for sending it's heartbeat
	client := internal.NewHeartbeatClient()

	go checkPeerHealth(registry)

	//For receiving side
	mux := http.NewServeMux()

	//HTTP request router, when a specific HTTP request arrives at this path, what function should handle this?
	mux.HandleFunc(
		"/v1/heartbeat",
		registry.HandleHeartbeat,
	)
	mux.HandleFunc(
		"/v1/status",
		registry.HandleStatus,
	)

	if *peerURL != "" {
		go sendHeartbeat(client, *peerURL, node)
	}

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

func sendHeartbeat(
	client *internal.HeartbeatClient,
	peerURL string,
	node internal.Node,
) {
	//ticker produces an event every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var sequence uint64
	//Occurs every second
	for range ticker.C {
		sequence++
		//Creates heartbeat describing THIS drone
		heartbeat := internal.Heartbeat{
			NodeID:   node.ID,
			Sequence: sequence,
			SentAt:   time.Now(),
			Position: node.Position,
			Battery:  node.Resources.Battery,
		}
		//Sends heartbeat to other drone
		if err := client.Send(peerURL, heartbeat); err != nil {
			log.Printf(
				"failed to send heartbeat to %s: %v",
				peerURL,
				err,
			)
		}
		continue
	}

	log.Printf(
		"sent heartbeat %d to %s",
		sequence,
		peerURL,
	)
}

func checkPeerHealth(registry *internal.PeerRegistry) {
	ticker := time.NewTicker(time.Second)

	defer ticker.Stop()

	for range ticker.C {
		registry.CheckHealth(time.Now())
	}

}
