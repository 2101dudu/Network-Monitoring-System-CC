package udp

import (
	"log"
	parse "nms/internal/jsonParse"
	utils "nms/internal/utils"
)

var agentsIPs map[byte][4]byte

func StartUDPServer(port string) {
	// incldude "| log.Lshortfile" in the log flags to include the file name and line of code in the log
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	jsonData := parse.GetDataFromJson("configs/tasks.json")
	var taskList []parse.Task = parse.ParseDataFromJson(jsonData)

	// validate tasks
	parse.ValidateTaskList(taskList)

	// Initialize the map
	agentsIPs = make(map[byte][4]byte)

	serverConn := utils.ResolveUDPAddrAndListen("localhost", "8081")
	handleRegistrations(serverConn)

	//serverConn.SetDeadline(time.Now().Add(5 * time.Second))

	//serverConn.Close()

	handleTasks(taskList)

	// go Receive metrics
}
