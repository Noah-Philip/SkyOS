package internal

import (
	"testing"
	"time"
)

func TestApplyHeartbeatCreatesReadyPeer(t *testing.T) {
	registry := NewPeerRegistry()
	now := time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)

	heartbeat := Heartbeat{
		NodeID:   "drone-2",
		Sequence: 1,
		Position: Position{
			X: 5,
			Y: 3,
			Z: 10,
		},
		Battery: 82,
	}

	registry.ApplyHeartbeat(heartbeat, now)

	peer, exists := registry.GetPeer("drone-2")
	if !exists {
		t.Fatal("expected drone-2 to exist in peer registry")
	}

	if peer.Status != PeerReady {
		t.Fatalf("expected status %q, got %q", PeerReady, peer.Status)
	}

	if peer.Battery != 82 {
		t.Fatalf("expected battery 82, got %v", peer.Battery)
	}

	if peer.Position != heartbeat.Position {
		t.Fatalf(
			"expected position %#v, got %#v",
			heartbeat.Position,
			peer.Position,
		)
	}

	if !peer.LastHeartbeatAt.Equal(now) {
		t.Fatalf(
			"expected received time %v, got %v",
			now,
			peer.LastHeartbeatAt,
		)
	}
}

func TestApplyHeartbeatIgnoresOlderSequence(t *testing.T) {
	registry := NewPeerRegistry()
	firstTime := time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)

	registry.ApplyHeartbeat(Heartbeat{
		NodeID:   "drone-2",
		Sequence: 2,
		Position: Position{X: 10, Y: 5, Z: 10},
		Battery:  80,
	}, firstTime)

	registry.ApplyHeartbeat(Heartbeat{
		NodeID:   "drone-2",
		Sequence: 1,
		Position: Position{X: 0, Y: 0, Z: 0},
		Battery:  10,
	}, firstTime.Add(time.Second))

	peer, exists := registry.GetPeer("drone-2")
	if !exists {
		t.Fatal("expected drone-2 to exist in peer registry")
	}

	if peer.LastSequence != 2 {
		t.Fatalf("expected sequence 2, got %d", peer.LastSequence)
	}

	if peer.Battery != 80 {
		t.Fatalf("expected battery 80, got %v", peer.Battery)
	}

	expectedPosition := Position{X: 10, Y: 5, Z: 10}
	if peer.Position != expectedPosition {
		t.Fatalf(
			"expected position %#v, got %#v",
			expectedPosition,
			peer.Position,
		)
	}

	if !peer.LastHeartbeatAt.Equal(firstTime) {
		t.Fatalf(
			"expected heartbeat time %v, got %v",
			firstTime,
			peer.LastHeartbeatAt,
		)
	}
}

func TestCheckHealthUpdatesPeerStatus(t *testing.T) {
	registry := NewPeerRegistry()
	now := time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC)

	registry.ApplyHeartbeat(Heartbeat{
		NodeID:   "drone-2",
		Sequence: 1,
		Battery:  80,
	}, now)

	registry.CheckHealth(now.Add(4 * time.Second))

	peer, _ := registry.GetPeer("drone-2")
	if peer.Status != PeerSuspect {
		t.Fatalf(
			"expected status %q, got %q",
			PeerSuspect,
			peer.Status,
		)
	}

	registry.CheckHealth(now.Add(7 * time.Second))

	peer, _ = registry.GetPeer("drone-2")
	if peer.Status != PeerDead {
		t.Fatalf(
			"expected status %q, got %q",
			PeerDead,
			peer.Status,
		)
	}
}
