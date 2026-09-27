package internal

import (
	"sync"
	"time"
)

type PeerStatus string

const (
	PeerReady   PeerStatus = "READY"
	PeerSuspect PeerStatus = "SUSPECT"
	PeerDead    PeerStatus = "DEAD"
)

const (
	PeerSuspectAfter = 3 * time.Second
	PeerDeadAfter    = 6 * time.Second
)

// Information a drone receives from another drone
type Heartbeat struct {
	NodeID   string
	Sequence uint64
	SentAt   time.Time
	Position Position
	Battery  float64
}

// Drone's locally stored knowledge about another drone
type Peer struct {
	ID              string
	Status          PeerStatus
	LastSequence    uint64
	LastHeartbeatAt time.Time
	Position        Position
	Battery         float64
}

// Drone's record of other drones it knows about
type PeerRegistry struct {
	mu sync.RWMutex
	//Key is drone ID, value is last known Peer state
	peers map[string]*Peer
}

func NewPeerRegistry() *PeerRegistry {
	return &PeerRegistry{
		peers: make(map[string]*Peer),
	}
}

// Gets another drone's information, store it in the registry if needed, otherwise change the information we know about it
func (r *PeerRegistry) ApplyHeartbeat(
	heartbeat Heartbeat,
	now time.Time,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	peer, exists := r.peers[heartbeat.NodeID]
	if !exists {
		peer = &Peer{
			ID: heartbeat.NodeID,
		}
		r.peers[heartbeat.NodeID] = peer
	}

	//Prevents old Heartbeat from applying.
	if heartbeat.Sequence <= peer.LastSequence {
		return
	}
	peer.LastSequence = heartbeat.Sequence
	peer.LastHeartbeatAt = now
	peer.Position = heartbeat.Position
	peer.Battery = heartbeat.Battery
	peer.Status = PeerReady
}

// Checks last heartbeat and determines if that drone is dead, suspected of being dead, or healthy.
func (r *PeerRegistry) CheckHealth(now time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, peer := range r.peers {
		timeSinceHeartbeat := now.Sub(peer.LastHeartbeatAt)

		switch {
		case timeSinceHeartbeat >= PeerDeadAfter:
			peer.Status = PeerDead
		case timeSinceHeartbeat >= PeerSuspectAfter:
			peer.Status = PeerSuspect

		default:
			peer.Status = PeerReady
		}
	}
}

func (r *PeerRegistry) GetPeer(id string) (Peer, bool) {
	//Read lock (multiple Goroutines can read at the same time)
	r.mu.RLock()
	defer r.mu.RUnlock()

	peer, exists := r.peers[id]

	if !exists {
		return Peer{}, false
	}

	//Copy of peer instead of the actual pointer
	return *peer, true
}

// Returns copies of each peer currently known by drone
func (r *PeerRegistry) ListPeers() []Peer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	peers := make([]Peer, 0, len(r.peers))

	for _, peer := range r.peers {
		peers = append(peers, *peer)
	}

	return peers
}
