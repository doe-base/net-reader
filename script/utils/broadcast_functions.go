package utils

import (
	"fmt"
	"net"
	"time"
)

func BroadcastListener() {
	// Listen on all interfaces at port 9999
	pc, err := net.ListenPacket("udp4", ":9999")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	fmt.Println("Listening for UDP broadcast on :9999...")

	buf := make([]byte, 1024)
	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}
		fmt.Printf("Received %d bytes from %s: %s\n", n, addr, string(buf[:n]))
	}
}
func BroadcastSender() {
	// Create a connection to send from
	// We use ":0" to let the OS pick any available local port
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Resolve the broadcast address
	dst, err := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")
	if err != nil {
		panic(err)
	}

	msg := []byte("Hello, everyone!")

	for {
		_, err := conn.WriteTo(msg, dst)
		if err != nil {
			fmt.Println("Send error:", err)
		} else {
			fmt.Println("Broadcast message sent.")
		}
		time.Sleep(2 * time.Second)
	}
}
