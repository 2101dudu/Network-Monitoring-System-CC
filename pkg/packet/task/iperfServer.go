package packet

// ------------------------- Iperf Server -----------------------------
type IperfServerMessage struct {
	SenderID  byte
	PacketID  byte
	TaskID    byte
	Frequency byte
	Command   string
} // ????????
