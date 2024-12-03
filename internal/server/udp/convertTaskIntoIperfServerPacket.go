package udp

import (
	parse "nms/internal/jsonParse"
	t "nms/internal/packet/task"
	"nms/internal/utils"
)

func ConvertTaskIntoIperfServerPacket(task parse.Task, serverIndex byte) t.IperfServerPacket {
	// main fields
	aID := task.Devices[serverIndex].DeviceID
	tID := task.TaskID
	freq := task.Frequency

	// build device metrics
	cpu := task.Devices[serverIndex].DeviceMetrics.CpuUsage
	ram := task.Devices[serverIndex].DeviceMetrics.RamUsage
	inter := task.Devices[serverIndex].DeviceMetrics.InterfaceStats
	devM := t.NewDeviceMetricsBuilder().SetCpuUsage(cpu).SetRamUsage(ram).SetInterfaceStats(inter).Build()

	// build alert flow conditions
	cpuUsage := task.Devices[serverIndex].AlertFlowConditions.CpuUsage
	ramUsage := task.Devices[serverIndex].AlertFlowConditions.RamUsage
	interStats := task.Devices[serverIndex].AlertFlowConditions.InterfaceStats
	packetLoss := task.Devices[serverIndex].AlertFlowConditions.PacketLoss
	jitter := task.Devices[serverIndex].AlertFlowConditions.Jitter
	alert := t.NewAlertFlowConditionsBuilder().SetCpuUsage(cpuUsage).SetRamUsage(ramUsage).SetInterfaceStats(interStats).SetPacketLoss(packetLoss).SetJitter(jitter).Build()

	// build iperf server command
	iperfCommand := "iperf -s -P 1"
	if !task.Devices[serverIndex].LinkMetrics.IperfParameters.Bandwidth {
		iperfCommand += " -u"
	}

	// build iperf server packet
	iperfPacket := t.NewIperfServerPacketBuilder().
		SetAgentID(aID).
		SetPacketID(utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)).
		SetTaskID(tID).
		SetFrequency(freq).
		SetDeviceMetrics(devM).
		SetAlertFlowConditions(alert).
		SetIperfServerCommand(iperfCommand).
		SetBandwidth(task.Devices[serverIndex].LinkMetrics.IperfParameters.Bandwidth).
		SetJitter(task.Devices[serverIndex].LinkMetrics.IperfParameters.Jitter).
		SetPacketLoss(task.Devices[serverIndex].LinkMetrics.IperfParameters.PacketLoss).
		Build()

	return iperfPacket
}
