package udp

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func getInterfaceStats(interfaceName string) (int, error) {
	// Read the content of ip command
	cmd := exec.Command("ip", "-s", "link", "show", interfaceName)
	data, err := cmd.CombinedOutput()

	if err != nil {
		return -1, err
	}

	lines := strings.Split(string(data), "\n")
	var receivedPackets, transmitedPackets int

	// Parse received packets
	receivedFields := strings.Fields(lines[3])
	if len(receivedFields) < 6 {
		log.Println("[ERROR 803] Unexpected ip command format")
	}

	receivedPackets, err = strconv.Atoi(receivedFields[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse received packets: %v", err)
	}
	if receivedPackets == 0 {
		return 0, fmt.Errorf("received zero packets, unexpected ip command format")
	}

	// Parse transmited packets
	transmitedFields := strings.Fields(lines[5])
	if len(transmitedFields) < 6 {
		log.Println("[ERROR 804] Unexpected /proc/meminfo format")
	}

	transmitedPackets, err = strconv.Atoi(transmitedFields[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse MemAvailable: %v", err)
	}

	// Calculate total packets
	totalPackets := receivedPackets + transmitedPackets
	return totalPackets, nil
}
