package task

import (
	"fmt"
	parse "nms/internal/jsonParse"
	"os"
)

func HandleTasks(taskList []parse.Task) {
	// for task in taskList

	for _, task := range taskList {
		// build, encode, and send ping
		if len(task.Devices) == 1 {
			// convert task into ping packet
			pingPacket := convertTaskIntoPingPacket(task)

			// encode ping packet
			data, err := EncodePingMessage(pingPacket)
			if err != nil {
				fmt.Println("[ERROR 21] Encoding ping packet")
				os.Exit(1)
			}

			// send ping packet
			// TODO

			// decode ping packet
			newPingPacket, err := DecodePingMessage(data[1:])
			if err != nil {
				fmt.Println("[ERROR 22] Decoding ping packet")
				os.Exit(1)
			}

			fmt.Print(newPingPacket.AgentID, newPingPacket.PacketID, newPingPacket.TaskID, newPingPacket.Frequency, newPingPacket.DeviceMetrics.CpuUsage, newPingPacket.DeviceMetrics.RamUsage, newPingPacket.DeviceMetrics.InterfaceStats, newPingPacket.AlertFlowConditions, newPingPacket.DeviceMetrics, newPingPacket.PingCommand+"\n\n")
		} else {
			if task.Devices[0].LinkMetrics.IperfParameters.IsServer {
				iperfServerPacket := convertTaskIntoIperfServerPacket(task, 0)
				iperfClientPacket := convertTaskIntoIperfClientPacket(task, 1)

				// encode iperf server packet
				dataServer, err := EncodeIperfServerPacket(iperfServerPacket)
				if err != nil {
					fmt.Println("[ERROR 25] Encoding iperf server packet")
					os.Exit(1)
				}

				// encode iperf client packet
				dataClient, err := EncodeIperfClientPacket(iperfClientPacket)
				if err != nil {
					fmt.Println("[ERROR 23] Encoding iperf client packet")
					os.Exit(1)
				}

				// send iperf client packet
				// TODO

				// send iperf server packet
				// TODO

				// decode iperf server packet
				newIperfServerPacket, err := DecodeIperfServerPacket(dataServer[1:])
				if err != nil {
					fmt.Println("[ERROR 26] Decoding iperf server packet")
					os.Exit(1)
				}

				// decode iperf client packet
				newIperfClientPacket, err := DecodeIperfClientPacket(dataClient[1:])
				if err != nil {
					fmt.Println("[ERROR 24] Decoding iperf client packet")
					os.Exit(1)
				}

				fmt.Print(newIperfServerPacket.AgentID, newIperfServerPacket.PacketID, newIperfServerPacket.TaskID, newIperfServerPacket.Frequency, newIperfServerPacket.DeviceMetrics.CpuUsage, newIperfServerPacket.DeviceMetrics.RamUsage, newIperfServerPacket.DeviceMetrics.InterfaceStats, newIperfServerPacket.AlertFlowConditions, newIperfServerPacket.DeviceMetrics, newIperfServerPacket.IperfServerCommand+"\n\n")
				fmt.Print(newIperfClientPacket.AgentID, newIperfClientPacket.PacketID, newIperfClientPacket.TaskID, newIperfClientPacket.Frequency, newIperfClientPacket.DeviceMetrics.CpuUsage, newIperfClientPacket.DeviceMetrics.RamUsage, newIperfClientPacket.DeviceMetrics.InterfaceStats, newIperfClientPacket.AlertFlowConditions, newIperfClientPacket.DeviceMetrics, newIperfClientPacket.IperfClientCommand+"\n\n")
			} else {
				iperfServerPacket := convertTaskIntoIperfServerPacket(task, 1)
				iperfClientPacket := convertTaskIntoIperfClientPacket(task, 0)

				// encode iperf server packet
				dataServer, err := EncodeIperfServerPacket(iperfServerPacket)
				if err != nil {
					fmt.Println("[ERROR 25] Encoding iperf server packet")
					os.Exit(1)
				}

				// encode iperf client packet
				dataClient, err := EncodeIperfClientPacket(iperfClientPacket)
				if err != nil {
					fmt.Println("[ERROR 23] Encoding iperf client packet")
					os.Exit(1)
				}

				// send iperf client packet
				// TODO

				// send iperf server packet
				// TODO

				// decode iperf server packet
				newIperfServerPacket, err := DecodeIperfServerPacket(dataServer[1:])
				if err != nil {
					fmt.Println("[ERROR 26] Decoding iperf server packet")
					os.Exit(1)
				}

				// decode iperf client packet
				newIperfClientPacket, err := DecodeIperfClientPacket(dataClient[1:])
				if err != nil {
					fmt.Println("[ERROR 24] Decoding iperf client packet")
					os.Exit(1)
				}

				fmt.Print(newIperfServerPacket.AgentID, newIperfServerPacket.PacketID, newIperfServerPacket.TaskID, newIperfServerPacket.Frequency, newIperfServerPacket.DeviceMetrics.CpuUsage, newIperfServerPacket.DeviceMetrics.RamUsage, newIperfServerPacket.DeviceMetrics.InterfaceStats, newIperfServerPacket.AlertFlowConditions, newIperfServerPacket.DeviceMetrics, newIperfServerPacket.IperfServerCommand+"\n\n")
				fmt.Print(newIperfClientPacket.AgentID, newIperfClientPacket.PacketID, newIperfClientPacket.TaskID, newIperfClientPacket.Frequency, newIperfClientPacket.DeviceMetrics.CpuUsage, newIperfClientPacket.DeviceMetrics.RamUsage, newIperfClientPacket.DeviceMetrics.InterfaceStats, newIperfClientPacket.AlertFlowConditions, newIperfClientPacket.DeviceMetrics, newIperfClientPacket.IperfClientCommand+"\n\n")
			}
		}

	}
}
