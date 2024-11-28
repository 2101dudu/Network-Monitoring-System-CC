package task

import (
	parse "nms/internal/jsonParse"
)

func ConvertTaskIntoIperfServerPacket(task parse.Task, serverIndex byte) IperfServerPacket {
	// main fields
	aID := task.Devices[serverIndex].DeviceID
	tID := task.TaskID
	freq := task.Frequency

	// build device metrics
	cpu := task.Devices[serverIndex].DeviceMetrics.CpuUsage
	ram := task.Devices[serverIndex].DeviceMetrics.RamUsage
	inter := task.Devices[serverIndex].DeviceMetrics.InterfaceStats
	devM := NewDeviceMetricsBuilder().SetCpuUsage(cpu).SetRamUsage(ram).SetInterfaceStats(inter).Build()

	// build alert flow conditions
	cpuUsage := task.Devices[serverIndex].AlertFlowConditions.CpuUsage
	ramUsage := task.Devices[serverIndex].AlertFlowConditions.RamUsage
	interStats := task.Devices[serverIndex].AlertFlowConditions.InterfaceStats
	packetLoss := task.Devices[serverIndex].AlertFlowConditions.PacketLoss
	jitter := task.Devices[serverIndex].AlertFlowConditions.Jitter
	alert := NewAlertFlowConditionsBuilder().SetCpuUsage(cpuUsage).SetRamUsage(ramUsage).SetInterfaceStats(interStats).SetPacketLoss(packetLoss).SetJitter(jitter).Build()

	// build iperf server command
	iperfCommand := "iperf3 -s"
	if !task.Devices[serverIndex].LinkMetrics.IperfParameters.Bandwidth {
		iperfCommand += " -u"
	}

	// build iperf server packet
	iperfPacket := NewIperfServerPacketBuilder().
		SetAgentID(aID).
		SetPacketID(1).
		SetTaskID(tID).
		SetFrequency(freq).
		SetDeviceMetrics(devM).
		SetAlertFlowConditions(alert).
		SetIperfServerCommand(iperfCommand).
		Build()

	return iperfPacket
}
