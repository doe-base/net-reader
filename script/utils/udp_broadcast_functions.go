package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type PeerInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

// Add a struct to track timestamp
type PeerStatus struct {
	Name     string
	LastSeen time.Time
}

/*
1️⃣ Listens on a UDP port: 0.0.0.0:9999
2️⃣ Broadcasts its identity every few seconds to: 255.255.255.255:9999 (every ip range). With Payload (example):

	{
		"name": "Daniel-PC",
		"ip": "192.168.0.161"
	}

3️⃣ Receives broadcasts from others and Adds them to peer list
*/
var Peers = make(map[string]PeerStatus)

func BroadcastListener() {
	pc, err := net.ListenPacket("udp4", ":9999")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	// Map to store discovered peers: [IP Address] -> Hostname
	myIP := getLocalIP()

	fmt.Println("Listening for peers...")

	buf := make([]byte, 1024)
	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}

		fmt.Printf("Packet physically received from: %s\n", addr.String())

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
		}

		// Update peer map with timestamp
		peerMu.Lock()
		if _, exists := Peers[info.IP]; !exists {
			fmt.Printf("✨ New Peer Discovered: %s at %s\n", info.Name, info.IP)
		}
		Peers[info.IP] = PeerStatus{
			Name:     info.Name,
			LastSeen: time.Now(),
		}
		peerMu.Unlock()
	}
}
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
		_, err := conn.Write(msg) // Now you can just use Write

		if err != nil {
			fmt.Println("Error:", err)
		}
		time.Sleep(5 * time.Second)
	}
}

// Helper to find the actual local network IP (not 127.0.0.1)
// RETURNS THE FIRST NON-LOOP IP
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

func CleanupPeers() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		peerMu.Lock()
		for ip, p := range Peers {
			if now.Sub(p.LastSeen) > 15*time.Second {
				fmt.Printf("⚠️ Peer expired: %s at %s\n", p.Name, ip)
				delete(Peers, ip)
			}
		}
		peerMu.Unlock()
	}
}
