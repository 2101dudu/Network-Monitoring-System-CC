package packet

// TYPE - TASK
type TaskType byte

const (
	IPERF TaskType = iota
	PING
	IPERFS
)

type DeviceMetrics struct {
	CpuUsage       bool
	RamUsage       bool
	InterfaceStats []string
}

type AlertFlowConditions struct {
	CpuUsage       byte
	RamUsage       byte
	InterfaceStats uint16
	PacketLoss     byte
	Jitter         uint16
}

/* func HandleTask(msgData []byte, senderID byte) error {
	msgTaskType := TaskType(msgData[0]) // Iperf/Ping
	msgPayload := msgData[1:]

	switch msgTaskType {
	case IPERFS:
		fmt.Println("[AGENT] Server Iperf Command received.")
		return nil
	case IPERF:
		fmt.Println("[AGENT] Iperf Command received.")

		task, err := DecodeIperfMessage(msgPayload)
		if err != nil {
			fmt.Println("[AGENT] [ERROR] Unable to decode Iperf message data:", err)
			return err
		}

		output, err := u.ExecuteCommand(task.Iperf.IperfCommand)

		return nil
	case PING:
		fmt.Println("[AGENT] Ping Command received.")

		task, err := DecodePingMessage(msgPayload)
		if err != nil {
			fmt.Println("[AGENT] [ERROR] Unable to decode Ping message data:", err)
			return err
		}

		output, err := u.ExecuteCommand(task.PingCommand)

		return nil
	default:
		err := errors.New("unknown task type in message data")
		fmt.Println("[AGENT] [ERROR]", err)
		return err
	}

} */
