package udp

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tcp "nms/internal/agent/tcp"
	alert "nms/internal/packet/alert"
	task "nms/internal/packet/task"
	utils "nms/internal/utils"
)

func ExecuteCommandWithMonitoring(command string, metrics task.DeviceMetrics, conditions task.AlertFlowConditions, taskID uint16) (string, error) {
	done := make(chan struct{})

	go monitorSystemMetrics(metrics, conditions, taskID, done)

	cmd := exec.Command("sh", "-c", command)
	//log.Println("Executing command:", command)
	stdout, err := cmd.CombinedOutput()

	close(done)

	if err != nil {
		return string(stdout), err
	}

	//log.Println("Command executed.")
	return string(stdout), nil
}

func monitorSystemMetrics(metrics task.DeviceMetrics, conditions task.AlertFlowConditions, taskID uint16, done chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	cpuAlert := false
	ramAlert := false

	//log.Println("[MONITOR] Start")
	for {

		if (ramAlert || !metrics.RamUsage) && (cpuAlert || !metrics.CpuUsage) { // If both alerts happened then can stop
			return
		}

		select {
		case <-done: // Command finished
			//log.Println("[MONITOR] End")
			return

		case <-ticker.C:
			// if there is not an alert of cpu and cpu Usage is a metric
			if !cpuAlert && metrics.CpuUsage {
				cpuUsage, errorCpu := getCpuUsage()
				if errorCpu != nil {
					fmt.Println("[AGENT] [MONITOR] [ERROR 180] Error getting CPU usage:", errorCpu)
					cpuUsage = 0
				}

				if cpuUsage > float32(conditions.CpuUsage) { // Build and send cpu alert
					cpuAlert = true

					alertTime := time.Now() // time of the alert

					newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
					buildAlert := alert.NewAlertBuilder().
						SetPacketID(newPacketID).
						SetSenderID(agentID).
						SetTaskID(taskID).
						SetAlertType(alert.CPU).
						SetExceeded(cpuUsage).
						SetTime(alertTime.Format("15:04:05.000000000"))

					newAlert := buildAlert.Build()                        // build full alert with given sets
					tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
					// fmt.Printf("[ALERT] CPU Usage exceeded: %.2f%% (limit: %d%%)\n", cpuUsage, conditions.CpuUsage)
				}
			}

			// if there is not an alert of ram and ram usage is a metric
			if !ramAlert && metrics.RamUsage {
				ramUsage, errorRam := getRamUsage()
				if errorRam != nil {
					fmt.Println("[AGENT] [MONITOR] [ERROR 181] Error getting RAM usage:", errorRam)
					ramUsage = 0
				}

				if ramUsage > float32(conditions.RamUsage) { // Build and send cpu alert
					ramAlert = true

					alertTime := time.Now() // time of the alert

					newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
					buildAlert := alert.NewAlertBuilder().
						SetPacketID(newPacketID).
						SetSenderID(agentID).
						SetTaskID(taskID).
						SetAlertType(alert.RAM).
						SetExceeded(ramUsage).
						SetTime(alertTime.Format("15:04:05.000000000"))

					newAlert := buildAlert.Build()                        // build full alert with given sets
					tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp
					// fmt.Printf("[ALERT] RAM Usage exceeded: %.2f%% (limit: %d%%)\n", ramUsage, conditions.RamUsage)
				}
			}
		}
	}
}

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
