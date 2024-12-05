package nettask

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getRamUsage() (float32, error) {
	// Read the content of /proc/meminfo
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, fmt.Errorf("failed to read /proc/meminfo: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var memTotal, memAvailable float64

	// Parse MemTotal
	memTotalFields := strings.Fields(lines[0])
	if len(memTotalFields) < 2 {
		log.Println("[ERROR 801] Unexpected /proc/meminfo format")
	}

	memTotal, err = strconv.ParseFloat(memTotalFields[1], 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse MemTotal: %v", err)
	}
	if memTotal == 0 {
		return 0, fmt.Errorf("MemTotal is zero, unexpected /proc/meminfo format")
	}

	// Parse MemAvailable
	memAvailableFields := strings.Fields(lines[2])
	if len(memAvailableFields) < 2 {
		log.Println("[ERROR 802] Unexpected /proc/meminfo format")
	}

	memAvailable, err = strconv.ParseFloat(memAvailableFields[1], 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse MemAvailable: %v", err)
	}

	// Calculate RAM usage percentage
	usage := 100 * (1 - memAvailable/memTotal)
	return float32(usage), nil
}
