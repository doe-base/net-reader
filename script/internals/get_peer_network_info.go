package internals

import (
	"net/http"
)

/*
		IP list can not be access just over any network
		To reliably discover all devices on a network, you typically need at least one of these:

		1. Control of the router
			-Access to DHCP tables
			-Access to ARP tables
			-Ability to disable client isolation

		2. Privileged network position
			-Being on the same subnet
	 		-No isolation
			-No firewall segmentation

		3. Network infrastructure access
			-Switch mirror ports (SPAN)
			-Managed switches
			-Enterprise monitoring tools


		If you are build a LAN where you have complete access from router config to firewall setting.
		Then there are tool in utils/lan_scanner_functions will help you scan and get every device on the network.
		This can be achieved via:
		1) ARP Scan      → Who exists on my LAN (most accurate)
		2) ICMP Ping     → Who responds to ping
		3) TCP Probe     → Who has open services
		4) mDNS / NBNS   → Who are you (hostname + device type)
		5) OS Fingerprint → What are you running
*/
func GetPeerNetworkInfo(w http.ResponseWriter, r http.Request) {

}
