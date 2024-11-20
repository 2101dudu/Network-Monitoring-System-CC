package task

type Iperf struct {
	TestDuration byte
	Bandwidth    bool
	Jitter       bool
	PacketLoss   bool
	IperfCommand string
}

type IperfBuilder struct {
	iperf Iperf
}

func NewIperfBuilder() *IperfBuilder {
	return &IperfBuilder{
		iperf: Iperf{
			TestDuration: 0,
			Bandwidth:    false,
			Jitter:       false,
			PacketLoss:   false,
			IperfCommand: "",
		},
	}
}

func (b *IperfBuilder) SetTestDuration(duration byte) *IperfBuilder {
	b.iperf.TestDuration = duration
	return b
}

func (b *IperfBuilder) SetBandwidth(enabled bool) *IperfBuilder {
	b.iperf.Bandwidth = enabled
	return b
}

func (b *IperfBuilder) SetJitter(enabled bool) *IperfBuilder {
	b.iperf.Jitter = enabled
	return b
}

func (b *IperfBuilder) SetPacketLoss(enabled bool) *IperfBuilder {
	b.iperf.PacketLoss = enabled
	return b
}

func (b *IperfBuilder) SetIperfCommand(command string) *IperfBuilder {
	b.iperf.IperfCommand = command
	return b
}

func (b *IperfBuilder) Build() Iperf {
	return b.iperf
}
