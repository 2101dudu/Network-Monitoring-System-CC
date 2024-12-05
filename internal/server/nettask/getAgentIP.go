package nettask

import (
	"fmt"
	"log"
	"nms/internal/utils"
)

func getAgentIP(agentID byte) string {
	exists := agentsIDs[agentID]
	if !exists {
		log.Fatalln(utils.Red+"[ERROR 35] Agent IP not found", utils.Reset)
	}
	return fmt.Sprintf("r%d", int(agentID))
}
