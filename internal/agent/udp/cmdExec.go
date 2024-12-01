package udp

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	task "nms/internal/packet/task"
)

func ExecuteCommandWithMonitoring(command string, metrics task.DeviceMetrics, conditions task.AlertFlowConditions) (string, bool, bool, error) {
	done := make(chan struct{})
	alertResults := make(chan struct {
		cpuAlert bool
		ramAlert bool
	})

	go monitorSystemMetrics(metrics, conditions, done, alertResults)

	cmd := exec.Command("sh", "-c", command)
	//log.Println("Executing command:", command)
	stdout, err := cmd.CombinedOutput()

	close(done)

	finalAlerts := <-alertResults

	if err != nil {
		return string(stdout), finalAlerts.cpuAlert, finalAlerts.ramAlert, err
	}

	/* 	if finalAlerts.cpuAlert {
	   		log.Println("alerta cpu.")
	   	}
	   	if finalAlerts.ramAlert {
	   		log.Println("alerta ram.")
	   	} */

	//log.Println("Command executed.")
	return string(stdout), finalAlerts.cpuAlert, finalAlerts.ramAlert, nil //os alertas poderemos retornar aqui para depois enviar um alerta geral com todas as coisas em que houve erro como mostrado na struct anterior
}

func monitorSystemMetrics(metrics task.DeviceMetrics, conditions task.AlertFlowConditions, done chan struct{}, alertResults chan struct {
	cpuAlert bool
	ramAlert bool
}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	cpuAlert := false
	ramAlert := false

	//log.Println("[MONITOR] Start")
	for {

		if ramAlert && cpuAlert { // If both alerts happened then can stop
			return
		}

		select {
		case <-done:

			alertResults <- struct {
				cpuAlert bool
				ramAlert bool
			}{cpuAlert, ramAlert}
			//log.Println("[MONITOR] End")
			return

		case <-ticker.C:
			cpuUsage, ramUsage := getCpuAndRamUsage()

			if !cpuAlert && metrics.CpuUsage && cpuUsage > float64(conditions.CpuUsage) {
				cpuAlert = true
				// fmt.Printf("[ALERT] CPU Usage exceeded: %.2f%% (limit: %d%%)\n", cpuUsage, conditions.CpuUsage)
			}
			if !ramAlert && metrics.RamUsage && ramUsage > float64(conditions.RamUsage) {
				ramAlert = true
				// fmt.Printf("[ALERT] RAM Usage exceeded: %.2f%% (limit: %d%%)\n", ramUsage, conditions.RamUsage)
			}
		}
	}
}

func getCpuAndRamUsage() (float64, float64) {
	// Uso da CPU
	cpuPercent, err := getCpuUsage()
	if err != nil {
		fmt.Println("[MONITOR] Error getting CPU usage:", err)
		cpuPercent = 0
	}

	// Uso de RAM
	memInfo, err := getRamUsage()
	if err != nil {
		fmt.Println("[MONITOR] Error getting RAM Usage:", err)
		return cpuPercent, 0
	}

	return cpuPercent, memInfo
}

func getCpuUsage() (float64, error) {
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
	usage := 100 * float64(total-idle) / float64(total)
	return usage, nil
}

func getRamUsage() (float64, error) {
	// Read the content of /proc/meminfo
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, fmt.Errorf("failed to read /proc/meminfo: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var memTotal, memAvailable float64

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		// Parse MemTotal
		if fields[0] == "MemTotal:" {
			memTotal, err = strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse MemTotal: %v", err)
			}
		}

		// Parse MemAvailable
		if fields[0] == "MemAvailable:" {
			memAvailable, err = strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse MemAvailable: %v", err)
			}
		}
	}

	if memTotal == 0 {
		return 0, fmt.Errorf("MemTotal is zero, unexpected /proc/meminfo format")
	}

	// Calculate RAM usage percentage
	usage := 100 * (1 - memAvailable/memTotal)
	return usage, nil
}
