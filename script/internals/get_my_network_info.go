package internals

import (
	"conceptual-lan/utils"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

func GetMyNetworkInfo(w http.ResponseWriter, r *http.Request) {

	// get all the deive interface (ethenet, wifi, vpn, usb tethering)
	ifaces, err := net.Interfaces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Filter for active interfaces that is not system loopback
	for index, iface := range ifaces {
		isUp := iface.Flags&net.FlagUp != 0
		isNotLoopback := iface.Flags&net.FlagLoopback == 0

		if isUp && isNotLoopback {
			ipv4, _ := filterForIPV4AndSubnet(iface.Addrs())
			ipAddress := ipv4.String()
			// Using Reverse DNS Lookup to get name of device name on network.
			// This only works if router or device has a DNS record
			names, err := net.LookupAddr(ipv4.IP.String())
			deviceName := "Unknown Device"
			if err == nil && len(names) > 0 {
				deviceName = names[0]
			}
			// using Local OUI Database to get vendor name
			vendor := utils.LookupOUI(iface.HardwareAddr)

			fmt.Fprintf(w, "Found Active Network %d: %s (MAC: %s) \n", index, iface.Name, iface.HardwareAddr)
			fmt.Fprintf(w, "Device Name: %s (IP: %s) (Vendor: %s)\n", deviceName, ipv4.IP.String(), vendor)
			fmt.Fprintf(w, "IP Address: %s\n", ipAddress)
			fmt.Fprintf(w, "\n")
		}
	}

	defaltGateway, _ := getDefaultGateway()
	fmt.Fprintf(w, "\nYour Machine's Default Gateway: %s\n", defaltGateway)

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
