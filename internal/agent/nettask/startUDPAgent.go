package nettask

import (
	utils "nms/internal/utils"
	"sync"
)

var (
	packetsWaitingAck = make(map[byte]bool)
	pMutex            sync.Mutex
	packetID          = byte(1)
	packetMutex       sync.Mutex
)

var agentID byte

func StartUDPAgent() {
	// get the ID of the agent
	agentID, _ = utils.GetAgentID()
	// make the IP of the agent
	agentIP := fmt.Sprintf("r%d", int(agentID))

	// make the agent open an UDP connection via port 9091
	agentConn := utils.ResolveUDPAddrAndListen(agentIP, "9091")

	// make the agent connect to the server via UDP on port 8081
	serverConn := utils.ResolveUDPAddrAndDial(utils.SERVERIP, "8081")

	// register the agent with the server
	registerAgent(serverConn)

	// handle tasks from the server
	handleTasks(agentConn)

	// close the agent connection
	agentConn.Close()
}
