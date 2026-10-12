package internal

import (
	"encoding/json"
	"net/http"
)

type LocalNodeStatus struct {
	ID         string     `json:"id"`
	Position   Position   `json:"position"`
	Battery    float64    `json:"battery"`
	NodeStatus NodeStatus `json:"status"`
}

type StatusResponse struct {
	Node  LocalNodeStatus `json:"node"`
	Peers []*Peer         `json:"peers"`
}

func NodeStatusHandler(
	node *Node,
	registry *PeerRegistry,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(
				w,
				"method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}
		response := StatusResponse{
			Node: LocalNodeStatus{
				ID:         node.ID,
				Position:   node.Position,
				Battery:    node.Resources.Battery,
				NodeStatus: node.Status,
			},
			Peers: registry.ListPeers(),
		}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(
				w,
				"failed to encode status response",
				http.StatusInternalServerError,
			)
		}
	}
}
