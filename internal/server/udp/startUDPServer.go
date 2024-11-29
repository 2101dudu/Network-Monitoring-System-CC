package udp

import (
	"log"
	parse "nms/internal/jsonParse"
	utils "nms/internal/utils"
	"sync"
)

var (
	packetsWaitingAck = make(map[byte]bool)
	pMutex            sync.Mutex
	packetID          = byte(1)
	packetMutex       sync.Mutex
)

var agentsIPs map[byte][4]byte

func StartUDPServer(port string) {
	// include "| log.Lshortfile" in the log flags to include the file name and line of code in the log
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	jsonData := parse.GetDataFromJson("configs/tasks.json")
	var taskList []parse.Task = parse.ParseDataFromJson(jsonData)

	// validate tasks
	parse.ValidateTaskList(taskList)

	// Initialize the map
	agentsIPs = make(map[byte][4]byte)

	// make the server open an UDP connection via port 8081
	serverConn := utils.ResolveUDPAddrAndListen("localhost", "8081")

	// handle registrations from agents
	handleRegistrations(serverConn)

	//serverConn.SetDeadline(time.Now().Add(5 * time.Second))

	// connect and send tasks to agents
	go handleTasks(taskList)

	// receive metrics from agents
	handleMetrics(serverConn)

	// close the server connection
	serverConn.Close()
}
