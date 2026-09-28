package internal

//Receiving side of heartbeat
//GET and POST are 2 different HTTP methods, they tell the server the kind of request the client is making

import (
	"encoding/json"
	"net/http"
	"time"
)

// Uses POST: Client sends data to server, server processes it and could change state.
func (r *PeerRegistry) HandleHeartbeat(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodPost {
		http.Error(
			writer,
			"only POST is allowed",
			http.StatusBadRequest,
		)
		return
	}
	defer request.Body.Close()

	var heartbeat Heartbeat

	if err := json.NewDecoder(request.Body).Decode(&heartbeat); err != nil {
		http.Error(
			writer,
			"invalid heartbeat JSON",
			http.StatusBadRequest,
		)
		return
	}

	if heartbeat.NodeID == "" {
		http.Error(
			writer,
			"heartbeat node_id REQUIRED",
			http.StatusBadRequest,
		)
		return
	}

	r.ApplyHeartbeat(heartbeat, time.Now())
	writer.WriteHeader(http.StatusNoContent)
}

// Uses GET: Client asks server for information, server reads state and responds.
func (r *PeerRegistry) HandleStatus(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		http.Error(
			writer,
			"only GET is allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	//Temporary response variable that holds the drone's peer list
	//Creates JSON text in the HTTP response & sends it to whoever made the request
	response := struct {
		Peers []Peer `json:"peers"`
	}{
		Peers: r.ListPeers(),
	}

	//Tells the HTTP client that the response body is JSON
	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	//Converts the response struct into JSON
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(
			writer,
			"failed to encode status",
			http.StatusInternalServerError,
		)
	}

}
