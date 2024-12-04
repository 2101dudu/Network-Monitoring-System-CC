package udp

import (
	"fmt"
	"log"
	"os/exec"
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
	stdout, err := cmd.CombinedOutput()

	close(done)

	if err != nil {
		return string(stdout), err
	}

	return string(stdout), nil
}

func monitorSystemMetrics(metrics task.DeviceMetrics, conditions task.AlertFlowConditions, taskID uint16, done chan struct{}) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	cpuHasExceeded := false
	ramHasExceeded := false

	for {
		if (ramHasExceeded || !metrics.RamUsage) && (cpuHasExceeded || !metrics.CpuUsage) && len(metrics.InterfaceStats) <= 0 { // If both alerts happened then can stop
			return
		}

		select {
		case <-done: // Command finished
			return

		case <-ticker.C:
			// if cpu usage has not been exceeded and is to be tracked
			go func() {
				if !cpuHasExceeded && metrics.CpuUsage {
					cpuHasExceeded = handleCpuUsage(conditions, taskID)
				}
			}()

			// if ram usage has not been exceeded and is to be tracked
			go func() {
				if !ramHasExceeded && metrics.RamUsage {
					ramHasExceeded = handleRamUsage(conditions, taskID)
				}
			}()

			go func() {
				if len(metrics.InterfaceStats) > 0 {
					for index, interfaceName := range metrics.InterfaceStats {
						go func() {
							packetsHaveExceeded := handleInterfaceStats(interfaceName, conditions, taskID)

							if packetsHaveExceeded {
								// Remove the interface from the list, as to not be checked again
								metrics.InterfaceStats = append(metrics.InterfaceStats[:index], metrics.InterfaceStats[index+1:]...)
							}
						}()
					}
				}
			}()
		}
	}
}

func handleCpuUsage(conditions task.AlertFlowConditions, taskID uint16) bool {
	cpuUsage, errorCpu := getCpuUsage()
	if errorCpu != nil {
		fmt.Println("[AGENT] [MONITOR] [ERROR 180] Error getting CPU usage:", errorCpu)
		cpuUsage = 0
		return false
	}

	// Build and send cpu alert
	if cpuUsage > float32(conditions.CpuUsage) {
		alertTime := time.Now() // time of the alert

		newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetSenderID(agentID).SetTaskID(taskID).SetAlertType(alert.CPU).SetExceeded(cpuUsage).SetTime(alertTime.Format("15:04:05.000000000"))

		newAlert := buildAlert.Build()                        // build full alert with given sets
		tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp

		return true
	}

	return false
}

func handleRamUsage(conditions task.AlertFlowConditions, taskID uint16) bool {
	ramUsage, errorRam := getRamUsage()
	if errorRam != nil {
		fmt.Println("[AGENT] [MONITOR] [ERROR 181] Error getting RAM usage:", errorRam)
		ramUsage = 0

		return false
	}

	// Build and send cpu alert
	if ramUsage > float32(conditions.RamUsage) {
		alertTime := time.Now() // time of the alert

		newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetSenderID(agentID).SetTaskID(taskID).SetAlertType(alert.RAM).SetExceeded(ramUsage).SetTime(alertTime.Format("15:04:05.000000000"))

		newAlert := buildAlert.Build()                        // build full alert with given sets
		tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp

		return true
	}

	return false
}

func handleInterfaceStats(interfaceName string, conditions task.AlertFlowConditions, taskID uint16) bool {
	interfaceStatsBefore, err := getInterfaceStats(interfaceName)
	if err != nil {
		log.Println("[ERROR 282] Error getting interface stats:", err)
		return false
	}
	time.Sleep(1 * time.Second)
	interfaceStatsAfter, err := getInterfaceStats(interfaceName)

	if err != nil {
		log.Println("[ERROR 182] Error getting interface stats:", err)
		return false
	}

	interfaceStats := interfaceStatsAfter - interfaceStatsBefore

	if interfaceStats > int(conditions.InterfaceStats) {
		alertTime := time.Now() // time of the alert

		newPacketID := utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)
		buildAlert := alert.NewAlertBuilder().SetPacketID(newPacketID).SetSenderID(agentID).SetTaskID(taskID).SetAlertType(alert.INTERFACESTATS).SetExceeded(float32(interfaceStats)).SetTime(alertTime.Format("15:04:05.000000000"))

		newAlert := buildAlert.Build()                        // build full alert with given sets
		tcp.ConnectTCPAndSendAlert(utils.SERVERTCP, newAlert) // Send an alert by tcp

		return true
	}

	return false
}
