package internals

import (
	"net"
	"net/http"
	"net/netip"
	"sync"
	"time"

	"github.com/mdlayher/arp"
)

// IPStatus represents the desired output format
type IPStatus struct {
	IP     string `json:"ip"`
	Active bool   `json:"active"`
}

func GetPeerNetworkInfo(w http.ResponseWriter, r http.Request) {

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

// CIDR → Iterate IPs → Send ARP → Collect replies
// ARP scanning tells you Which devices physically exist on my network
func ARPScan(ipnet *net.IPNet, iface *net.Interface) (map[string]string, error) {
	results := make(map[string]string)

	c, err := arp.Dial(iface)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	ip := ipnet.IP.Mask(ipnet.Mask)

	for ip := append(net.IP(nil), ip...); ipnet.Contains(ip); incrementIP(ip) {
		if ip.Equal(ipnet.IP) {
			continue
		}

		err := c.SetDeadline(time.Now().Add(300 * time.Millisecond))
		if err != nil {
			continue
		}

		addr, ok := netip.AddrFromSlice(ip)
		if !ok {
			continue
		}

		hw, err := c.Resolve(addr)
		if err == nil {
			results[ip.String()] = hw.String()
		}
	}

	return results, nil
}
