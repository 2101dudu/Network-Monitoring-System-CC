package utils

import (
	"os/exec"
)

type MessageType byte

const (
	ACK              MessageType = iota // iota = 0
	REGISTRATION                        // iota = 1
	METRICSGATHERING                    // iota = 2
	TASK
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
