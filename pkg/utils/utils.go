package utils

type MessageType byte

const (
	ACK              MessageType = iota // iota = 0
	REGSITRATION                        // iota = 1
	METRICSGATHERING                    // iota = 2
	ERROR                               // iota = 3
)

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
