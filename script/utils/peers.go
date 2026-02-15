package utils

import "sync"

type DiscoveredPeer struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

var (
	peerMu sync.RWMutex
	peers  = make(map[string]DiscoveredPeer)
)

func AddPeer(p DiscoveredPeer) {
	peerMu.Lock()
	defer peerMu.Unlock()
	peers[p.IP] = p
}

func GetPeers() []DiscoveredPeer {
	peerMu.RLock()
	defer peerMu.RUnlock()

	out := make([]DiscoveredPeer, 0, len(peers))
	for _, p := range peers {
		out = append(out, p)
	}
	return out
}
