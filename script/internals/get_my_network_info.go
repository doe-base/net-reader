package internals

import (
	"conceptual-lan/utils"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

type NetworkDevice struct {
	Index      int    `json:"index"`
	Interface  string `json:"interface"`
	MAC        string `json:"mac"`
	DeviceName string `json:"deviceName"`
	IP         string `json:"ip"`
	Vendor     string `json:"vendor"`
	LocalIP    string `json:"localIp"` // <-- json tag
}

func GetMyNetworkInfo(w http.ResponseWriter, r *http.Request) {
	// get all the device interfaces (ethernet, wifi, vpn, usb tethering)
	ifaces, err := net.Interfaces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	devices := []NetworkDevice{}

	// Filter for active interfaces that are not loopback
	for index, iface := range ifaces {
		isUp := iface.Flags&net.FlagUp != 0
		isNotLoopback := iface.Flags&net.FlagLoopback == 0

		if isUp && isNotLoopback {
			ipv4, _ := filterForIPV4AndSubnet(iface.Addrs())
			if ipv4 == nil {
				continue
			}

			ipAddress := ipv4.String()

			// Reverse DNS lookup for device name
			deviceName := "Unknown Device"
			names, err := net.LookupAddr(ipv4.IP.String())
			if err == nil && len(names) > 0 {
				deviceName = names[0]
			}

			// Vendor lookup
			vendor := utils.LookupOUI(iface.HardwareAddr)

			device := NetworkDevice{
				Index:      index,
				Interface:  iface.Name,
				MAC:        iface.HardwareAddr.String(),
				DeviceName: deviceName,
				IP:         ipv4.IP.String(),
				Vendor:     vendor,
				LocalIP:    ipAddress,
			}

			devices = append(devices, device)
		}
	}

	// --- 2️⃣ Get default gateway ---
	defGateway, _ := getDefaultGateway()

	// --- 3️⃣ Include discovered peers ---
	// For now, we assume a global peers map updated by BroadcastListener
	// You can make it concurrent safe with a mutex if needed
	peersList := []utils.PeerInfo{}
	for ip, status := range utils.Peers {
		peersList = append(peersList, utils.PeerInfo{
			Name: status.Name,
			IP:   ip,
		})
	}

	// --- 4️⃣ Respond with combined JSON ---
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"count":          len(devices),
		"devices":        devices,
		"defaultGateway": defGateway,
		"peers":          peersList,
	})
}

// one iface.Addrs() can return multiple address object (of type: *net.IPNet, *net.IPAddr etc)
func filterForIPV4AndSubnet(addrs []net.Addr, err error) (*net.IPNet, error) {
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			// Check if the IP inside the network is IPv4
			if v.IP.To4() != nil {
				return v, nil
			}
		}
	}

	return nil, fmt.Errorf("no IPv4 address found")
}

func getDefaultGateway() (string, error) {
	// 1. Run the command: ip route show default
	out, err := exec.Command("ip", "route", "show", "default").Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	// 2. Parse the output (Example: "default via 192.168.1.1 dev eth0")
	fields := strings.Fields(string(out))
	if len(fields) >= 3 && fields[1] == "via" {
		return fields[2], nil
	}

	return "", fmt.Errorf("could not find gateway in output: %s", string(out))
}
