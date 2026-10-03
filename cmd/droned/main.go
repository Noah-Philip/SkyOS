package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
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

	configPath := flag.String(
		"config",
		"",
		"path to the swarm config file",
	)

	//Reads in arguments passed in command line & fills in NodeID
	flag.Parse()

	if *configPath == "" {
		log.Fatal("--config is required")
	}

	config, err := internal.LoadSwarmConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	nodeConfig, found := config.NodeByID(*NodeID)
	if !found {
		log.Fatalf("node %q was not found in config", *NodeID)
	}

	node := internal.Node{
		ID: nodeConfig.ID,
		Resources: internal.Resources{
			Battery: nodeConfig.Battery,
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

	if len(nodeConfig.Peers) > 0 {
		go sendHeartbeat(client, nodeConfig.Peers, node)
	}

	address := fmt.Sprintf(":%d", nodeConfig.Port)

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
	peers []string,
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
		for _, peerURL := range peers {
			//Sends heartbeat to other drone
			if err := client.Send(peerURL, heartbeat); err != nil {
				log.Printf(
					"failed to send heartbeat to %s: %v",
					peerURL,
					err,
				)
				continue
			}
			log.Printf(
				"sent heartbeat %d to %s",
				sequence,
				peerURL,
			)
		}
	}
}

func checkPeerHealth(registry *internal.PeerRegistry) {
	ticker := time.NewTicker(time.Second)

	defer ticker.Stop()

	for range ticker.C {
		registry.CheckHealth(time.Now())
	}

}

type peerURLs []string

func (p *peerURLs) String() string {
	return strings.Join(*p, ",")
}
func (p *peerURLs) Set(value string) error {
	*p = append(*p, value)
	return nil
}
