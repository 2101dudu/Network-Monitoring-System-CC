package udp

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	task "nms/internal/packet/task"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func ExecuteCommandWithMonitoring(command string, metrics task.DeviceMetrics, conditions task.AlertFlowConditions) (string, bool, bool, error) {
	done := make(chan struct{})
	alertResults := make(chan struct {
		cpuAlert bool
		ramAlert bool
	})

	go monitorSystemMetrics(metrics, conditions, done, alertResults)

	cmd := exec.Command("sh", "-c", command)
	log.Println("Executing command:", command)
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

	log.Println("Command executed.")
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

	log.Println("[MONITOR] Start")
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
			log.Println("[MONITOR] End")
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
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println("[MONITOR] Erro ao coletar CPU Usage:", err)
		cpuPercent = []float64{0}
	}

	// Uso de RAM
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("[MONITOR] Erro ao coletar RAM Usage:", err)
		return 0, 0
	}

	return cpuPercent[0], memInfo.UsedPercent
}
