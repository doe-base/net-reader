package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func DownloadOUI() {
	url := "https://standards-oui.ieee.org/oui/oui.txt"
	fileName := "oui.txt"

	log.Printf("Downloading OUI database from %s...", url)

	// 1. Create a custom client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// 2. Add a User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// 3. Execute the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error connecting to IEEE:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server returned error: %s", resp.Status)
	}

	// 4. Create the local file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	// 5. Stream the data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal("Error saving data:", err)
	}

	log.Println("Download complete! Saved as oui.txt")
}

// ParseOUI reads the downloaded file and returns a map[Prefix]VendorName
func ParseOUI(filename string) (map[string]string, error) {
	vendors := make(map[string]string)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// We only care about the lines that have the (hex) label
		if strings.Contains(line, "(hex)") {
			// 1. Split by the "(hex)" marker
			parts := strings.Split(line, "(hex)")
			if len(parts) < 2 {
				continue
			}

			// 2. Clean the Prefix (Left side)
			// This takes "28-6F-B9   " and makes it "28-6F-B9"
			prefix := strings.TrimSpace(parts[0])

			// 3. Clean the Vendor (Right side)
			// This handles the tabs and spaces after (hex)
			vendor := strings.TrimSpace(parts[1])

			vendors[prefix] = vendor
		}
	}
	return vendors, scanner.Err()
}

func LookupOUI(macAddress net.HardwareAddr) string {
	// 1. Check if the interface actually has a MAC address (loopback interfaces don't)
	if len(macAddress) >= 3 {

		// 2. Format the first 3 bytes as a hex string (e.g., "00e04c")
		// %x prints the hex; [:3] takes the first 3 bytes (the OUI)
		oui := fmt.Sprintf("%02X-%02X-%02X",
			macAddress[0],
			macAddress[1],
			macAddress[2],
		)

		// 3. Look it up in your map
		vendorMap, _ := ParseOUI("oui.txt")
		vendor := vendorMap[oui]

		if vendor == "" {
			vendor = "Unknown Vendor"
		}

		return vendor
	}
	return "Unknown Vendor"
}
