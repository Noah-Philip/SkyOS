package internal

//Receiving side of heartbeat

import (
	"encoding/json"
	"net/http"
	"time"
)

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
