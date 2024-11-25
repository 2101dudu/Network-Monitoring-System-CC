package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type PacketType byte

const (
	ACK              PacketType = iota // iota = 0
	REGISTRATION                       // iota = 1
	METRICSGATHERING                   // iota = 2
	IPERFCLIENT
	IPERFSERVER
	PING
)

const (
	TIMEOUTSECONDS = 2
	MAXAGENTS      = 1
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

func IPStringToByte(ip string) ([4]byte, error) {
	var byteIP [4]byte
	segments := strings.Split(ip, ".")
	if len(segments) != 4 {
		return byteIP, fmt.Errorf("invalid IP address: %s", ip)
	}

	for i, segment := range segments {
		intIP, err := strconv.Atoi(segment)
		if err != nil || intIP < 0 || intIP > 255 {
			return byteIP, fmt.Errorf("invalid IP segment: %s", segment)
		}
		byteIP[i] = byte(intIP)
	}

	return byteIP, nil
}
