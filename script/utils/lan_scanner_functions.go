package utils

import (
	"fmt"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/mdlayher/arp"
)

type IPStatus struct {
	IP     string `json:"ip"`
	Active bool   `json:"active"`
}

// Use the CIDR to get the ip range
func GenerateIPList(ipnet *net.IPNet) []IPStatus {
	var ips []IPStatus

	ip := ipnet.IP.To4()
	if ip == nil {
		return nil // not IPv4
	}

	ip = ip.Mask(ipnet.Mask)

	for ipnet.Contains(ip) {
		ips = append(ips, IPStatus{
			IP:     ip.String(),
			Active: false,
		})

		for i := len(ip) - 1; i >= 0; i-- {
			ip[i]++
			if ip[i] != 0 {
				break
			}
		}
	}

	// remove network + broadcast
	if len(ips) > 2 {
		return ips[1 : len(ips)-1]
	}
	return nil
}

/*
	 APR Scanning
		CIDR → Iterate IPs → Send ARP → Collect replies
		ARP scanning tells you Which devices physically exist on my network
*/

func ARPScan(ipnet *net.IPNet, iface *net.Interface) ([]IPStatus, error) {
	ips := GenerateIPList(ipnet)
	if len(ips) == 0 {
		return nil, fmt.Errorf("no IPs generated")
	}

	c, err := arp.Dial(iface)
	if err != nil {
		fmt.Println("ARP DIAL ERROR:", err)
		return nil, err
	}
	defer c.Close()

	for i := range ips {
		ip := net.ParseIP(ips[i].IP)
		if ip == nil {
			continue
		}

		addr, ok := netip.AddrFromSlice(ip)
		if !ok {
			continue
		}

		_ = c.SetDeadline(time.Now().Add(250 * time.Millisecond))

		if _, err := c.Resolve(addr); err == nil {
			ips[i].Active = true
		}
	}

	return ips, nil
}

/*
	 TCP port scan method
		Using Goroutines and a WaitGroup to scan all IPs simultaneously
		Scanning on multiple ports to get best result.
		Using Worker Pool to prevent goroutine explosion
*/
func TCPNetworkScan(ipnet *net.IPNet, workers int) []string {
	var wg sync.WaitGroup

	ipsChan := make(chan string, workers)
	resultsChan := make(chan string, workers)

	ports := []string{"80", "443", "22", "445", "139"}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for targetIP := range ipsChan {
				for _, port := range ports {
					conn, err := net.DialTimeout("tcp", targetIP+":"+port, 300*time.Millisecond)
					if err == nil {
						conn.Close()
						resultsChan <- targetIP
						break
					}
				}
			}
		}()
	}

	ip := ipnet.IP.To4()
	if ip == nil {
		return nil
	}

	startIP := ip.Mask(ipnet.Mask)
	for ip = append(net.IP(nil), startIP...); ipnet.Contains(ip); incrementIP(ip) {
		ipsChan <- ip.String()
	}

	close(ipsChan)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect only ACTIVE IPs
	var active []string
	for ip := range resultsChan {
		active = append(active, ip)
	}

	return active
}

// Helper to increment IP bytes
func incrementIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}
