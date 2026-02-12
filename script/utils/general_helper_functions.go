package utils

import "net"

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
