package utils

import (
	"os/exec"
	s "strconv"
)

type MessageType byte

const (
	ACK              MessageType = iota // iota = 0
	REGISTRATION                        // iota = 1
	METRICSGATHERING                    // iota = 2
	TASK
)

const (
	TIMEOUTSECONDS = 2
	MAXAGENTS      = 10
)

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func GetAgentID() (byte, error) {
	cmd := exec.Command("whoami")
	_, err := cmd.Output()

	// requires string parsing ignoring all characters (e.g.: "PC1" -> 1; "router2" -> 2)
	//id := whoami
	return byte(1), err
}

func IPStringToByte(ip string) [4]byte {
	byteIP := [4]byte{0, 0, 0, 0}
	n := 0
	for i, aux := 0, 0; i < len(ip); i++ {
		if ip[i] == '.' {
			intIP, _ := s.Atoi(ip[aux:i])
			byteIP[n] = byte(intIP)
			n++
			aux = i + 1
		}
	}
	return byteIP
}
