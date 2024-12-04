package udp

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getCpuUsage() (float32, error) {
	// Read the content of /proc/stat
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, fmt.Errorf("failed to read /proc/stat: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	cpuLine := lines[0]
	fields := strings.Fields(cpuLine)

	// Ensure we have at least the first 8 fields
	if len(fields) < 8 {
		return 0, fmt.Errorf("unexpected format in /proc/stat")
	}

	// Convert fields to integers
	var values []int64
	for _, field := range fields[1:] {
		value, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse CPU value: %v", err)
		}
		values = append(values, value)
	}

	// Calculate total and idle time
	total := values[0] + values[1] + values[2] + values[3] + values[4] + values[5] + values[6]
	idle := values[3]

	// Calculate CPU usage percentage
	usage := 100 * float32(total-idle) / float32(total)
	return usage, nil
}
