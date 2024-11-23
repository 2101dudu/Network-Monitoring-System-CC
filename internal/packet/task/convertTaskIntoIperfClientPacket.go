package task

import (
	"fmt"
	parse "nms/internal/jsonParse"
)

func convertTaskIntoIperfClientPacket(task parse.Task, clientIndex byte) IperfClientPacket {
	// main fields
	aID := task.Devices[clientIndex].DeviceID
	tID := task.TaskID
	freq := task.Frequency

	// build device metrics
	cpu := task.Devices[clientIndex].DeviceMetrics.CpuUsage
	ram := task.Devices[clientIndex].DeviceMetrics.RamUsage
	inter := task.Devices[clientIndex].DeviceMetrics.InterfaceStats
	devM := NewDeviceMetricsBuilder().SetCpuUsage(cpu).SetRamUsage(ram).SetInterfaceStats(inter).Build()

	// build alert flow conditions
	cpuUsage := task.Devices[clientIndex].AlertFlowConditions.CpuUsage
	ramUsage := task.Devices[clientIndex].AlertFlowConditions.RamUsage
	interStats := task.Devices[clientIndex].AlertFlowConditions.InterfaceStats
	packetLoss := task.Devices[clientIndex].AlertFlowConditions.PacketLoss
	jitter := task.Devices[clientIndex].AlertFlowConditions.Jitter
	alert := NewAlertFlowConditionsBuilder().SetCpuUsage(cpuUsage).SetRamUsage(ramUsage).SetInterfaceStats(interStats).SetPacketLoss(packetLoss).SetJitter(jitter).Build()

	// build iperf client command
	serverIndex := 1 - clientIndex
	destination := task.Devices[serverIndex].DeviceID

	iperfCommand := fmt.Sprintf("iperf3 -c %d", destination)
	if !task.Devices[clientIndex].LinkMetrics.IperfParameters.Bandwidth {
		iperfCommand += " -u"
	}

	// build iperf client packet
	iperfPacket := NewIperfClientPacketBuilder().
		SetAgentID(aID).
		SetPacketID(1).
		SetTaskID(tID).
		SetFrequency(freq).
		SetDeviceMetrics(devM).
		SetAlertFlowConditions(alert).
		SetIperfClientCommand(iperfCommand).
		SetBandwidth(task.Devices[0].LinkMetrics.IperfParameters.Bandwidth).
		SetJitter(task.Devices[0].LinkMetrics.IperfParameters.Jitter).
		SetPacketLoss(task.Devices[0].LinkMetrics.IperfParameters.PacketLoss).
		Build()

	return iperfPacket
}
