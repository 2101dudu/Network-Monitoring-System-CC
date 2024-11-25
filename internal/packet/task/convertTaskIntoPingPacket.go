package task

import (
	"fmt"
	parse "nms/internal/jsonParse"
)

func convertTaskIntoPingPacket(task parse.Task) PingPacket {
	// main fields
	aID := task.Devices[0].DeviceID
	tID := task.TaskID
	freq := task.Frequency

	// build device metrics
	cpu := task.Devices[0].DeviceMetrics.CpuUsage
	ram := task.Devices[0].DeviceMetrics.RamUsage
	inter := task.Devices[0].DeviceMetrics.InterfaceStats
	devM := NewDeviceMetricsBuilder().SetCpuUsage(cpu).SetRamUsage(ram).SetInterfaceStats(inter).Build()

	// build alert flow conditions
	cpuUsage := task.Devices[0].AlertFlowConditions.CpuUsage
	ramUsage := task.Devices[0].AlertFlowConditions.RamUsage
	interStats := task.Devices[0].AlertFlowConditions.InterfaceStats
	packetLoss := task.Devices[0].AlertFlowConditions.PacketLoss
	jitter := task.Devices[0].AlertFlowConditions.Jitter

	alert := NewAlertFlowConditionsBuilder().SetCpuUsage(cpuUsage).SetRamUsage(ramUsage).SetInterfaceStats(interStats).SetPacketLoss(packetLoss).SetJitter(jitter).Build()

	// build ping command
	destination := task.Devices[0].LinkMetrics.PingParameters.Destination
	pingCommand := fmt.Sprintf("ping -c %d -i %d %s", task.Devices[0].LinkMetrics.PingParameters.PacketCount, task.Devices[0].LinkMetrics.PingParameters.Frequency, destination)

	// build ping packet
	pingPacket := NewPingPacketBuilder().
		SetAgentID(aID).
		SetPacketID(1).
		SetTaskID(tID).
		SetFrequency(freq).
		SetDeviceMetrics(devM).
		SetAlertFlowConditions(alert).
		SetPingCommand(pingCommand).
		Build()

	return pingPacket
}
