package discovery

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type PeerInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}
type PeerStatus struct {
	Name     string
	LastSeen time.Time
}

var (
	PeerMu sync.RWMutex
	Peers  = make(map[string]PeerStatus)
)

// Listens on :9999, receives peer broadcasts, ignores its own messages,
// and stores discovered peers in a shared map with a LastSeen timestamp.
func BroadcastListener() {
	pc, err := net.ListenPacket("udp4", ":9999")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	myIP := getLocalIP()

	fmt.Println("Listening for peers...")

	buf := make([]byte, 1024)
	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}

		// 1. Decode the JSON
		var info PeerInfo
		err = json.Unmarshal(buf[:n], &info)
		if err != nil {
			fmt.Println("Received invalid data")
			continue
		}

		// Ignore self
		if info.IP == myIP {
			continue
		} else {
			fmt.Printf("Packet physically received from: %s\n", addr.String())
		}

		// Update peer map with timestamp
		PeerMu.Lock()
		if _, exists := Peers[info.IP]; !exists {
			fmt.Printf("✨ New Peer Discovered: %s at %s\n", info.Name, info.IP)
		}
		Peers[info.IP] = PeerStatus{
			Name:     info.Name,
			LastSeen: time.Now(),
		}
		PeerMu.Unlock()
	}
}

// Every 5 seconds, it broadcasts this machine’s hostname and local IP as JSON to 255.255.255.255:9999.
func BroadcastSender() {
	hostname, _ := os.Hostname()

	// Use DialUDP to setup a connection to the broadcast address
	dst, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")
	conn, err := net.DialUDP("udp4", nil, dst)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		info := PeerInfo{
			Name: hostname,
			IP:   getLocalIP(),
		}
		msg, _ := json.Marshal(info)
		_, err := conn.Write(msg)

		if err != nil {
			fmt.Println("Error:", err)
		}
		time.Sleep(5 * time.Second)
	}
}

// Finds the first non-loopback IPv4 address of the machine.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// Runs periodically and removes peers that haven’t been seen for 15 seconds.
func CleanupPeers() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		PeerMu.Lock()
		for ip, p := range Peers {
			if now.Sub(p.LastSeen) > 15*time.Second {
				fmt.Printf("⚠️ Peer expired: %s at %s\n", p.Name, ip)
				delete(Peers, ip)
			}
		}
		PeerMu.Unlock()
	}
}
