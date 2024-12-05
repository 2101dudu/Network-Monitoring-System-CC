package nettask

import (
	"fmt"
	parse "nms/internal/jsonParse"
	t "nms/internal/packet/task"
	"nms/internal/utils"
)

func ConvertTaskIntoIperfClientPacket(task parse.Task, clientIndex byte) t.IperfClientPacket {
	// main fields
	aID := task.Devices[clientIndex].DeviceID
	tID := task.TaskID
	freq := task.Frequency

	// build device metrics
	cpu := task.Devices[clientIndex].DeviceMetrics.CpuUsage
	ram := task.Devices[clientIndex].DeviceMetrics.RamUsage
	inter := task.Devices[clientIndex].DeviceMetrics.InterfaceStats
	devM := t.NewDeviceMetricsBuilder().SetCpuUsage(cpu).SetRamUsage(ram).SetInterfaceStats(inter).Build()

	// build alert flow conditions
	cpuUsage := task.Devices[clientIndex].AlertFlowConditions.CpuUsage
	ramUsage := task.Devices[clientIndex].AlertFlowConditions.RamUsage
	interStats := task.Devices[clientIndex].AlertFlowConditions.InterfaceStats
	packetLoss := task.Devices[clientIndex].AlertFlowConditions.PacketLoss
	jitter := task.Devices[clientIndex].AlertFlowConditions.Jitter
	alert := t.NewAlertFlowConditionsBuilder().SetCpuUsage(cpuUsage).SetRamUsage(ramUsage).SetInterfaceStats(interStats).SetPacketLoss(packetLoss).SetJitter(jitter).Build()

	// build iperf client command
	serverIndex := 1 - clientIndex
	destinationID := task.Devices[serverIndex].DeviceID
	testDuration := task.Devices[clientIndex].LinkMetrics.IperfParameters.TestDuration

	iperfCommand := fmt.Sprintf("iperf -c r%d -t %d", destinationID, testDuration)
	if !task.Devices[clientIndex].LinkMetrics.IperfParameters.Bandwidth {
		iperfCommand += " -u"
	}

	// build iperf client packet
	iperfPacket := t.NewIperfClientPacketBuilder().
		SetAgentID(aID).
		SetPacketID(utils.ReadAndIncrementPacketID(&packetID, &packetMutex, true)).
		SetTaskID(tID).
		SetFrequency(freq).
		SetDeviceMetrics(devM).
		SetAlertFlowConditions(alert).
		SetIperfClientCommand(iperfCommand).
		SetBandwidth(task.Devices[clientIndex].LinkMetrics.IperfParameters.Bandwidth).
		SetJitter(task.Devices[clientIndex].LinkMetrics.IperfParameters.Jitter).
		SetPacketLoss(task.Devices[clientIndex].LinkMetrics.IperfParameters.PacketLoss).
		Build()

	return iperfPacket
}
