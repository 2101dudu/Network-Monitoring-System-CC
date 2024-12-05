package nettask

import (
	"fmt"
	"log"
)

func getAgentIP(deviceID byte) string {
	agentIPByte, exists := agentsIPs[deviceID]
	if !exists {
		log.Fatalln("[ERROR 35] Agent IP not found")
	}
	return fmt.Sprintf("%d.%d.%d.%d", agentIPByte[0], agentIPByte[1], agentIPByte[2], agentIPByte[3])
}
