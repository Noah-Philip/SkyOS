package internal

//Sending side of heartbeat!

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// http client: starts requests
type HeartbeatClient struct {
	client *http.Client
}

func NewHeartbeatClient() *HeartbeatClient {
	return &HeartbeatClient{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (c *HeartbeatClient) Send(
	peerBaseURL string,
	heartbeat Heartbeat,
) error {
	//Converts the heartbeat to JSON
	data, err := json.Marshal(heartbeat)
	if err != nil {
		return fmt.Errorf("encode heartbeat: %w", err)
	}

	endpoint := strings.TrimRight(
		peerBaseURL,
		"/",
	) + "/v1/heartbeat"

	//http POST request
	request, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("encode heartbeat: %w", err)
	}

	//labels it as JSON
	request.Header.Set(
		"Content-Type",
		"application/json",
	)
	//sends request
	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf(
			"sent heartbeat at %q: %w",
			peerBaseURL,
			err,
		)
	}
	defer response.Body.Close()

	//Verifies it works, receiving handler returns 204 no content for now
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf(
			"unexpected heartbeat response: %s",
			response.Status,
		)
	}
	return nil

}
