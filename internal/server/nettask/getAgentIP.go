package nettask

import (
	"fmt"
	"log"
)

func getAgentIP(agentID byte) string {
	exists := agentsIDs[agentID]
	if !exists {
		log.Fatalln("[ERROR 35] Agent IP not found")
	}
	return fmt.Sprintf("r%d", int(agentID))
}
