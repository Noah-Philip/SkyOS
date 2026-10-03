package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

type NodeConfig struct {
	ID      string   `json: "id"`
	Port    int      `json: "port"`
	Battery float64  `json: "battery"`
	Peers   []string `json: "peers"`
}

type SwarmConfig struct {
	Nodes []NodeConfig `json:"nodes"`
}

func LoadSwarmConfig(path string) (*SwarmConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read swarm config: %w", err)
	}
	var config SwarmConfig

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse swarm config: %w", err)
	}

	return &config, nil
}

func (c *SwarmConfig) NodeByID(id string) (NodeConfig, bool) {
	for _, node := range c.Nodes {
		if node.ID == id {
			return node, true
		}
	}

	return NodeConfig{}, false
}
