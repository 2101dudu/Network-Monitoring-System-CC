package udp

import (
	"nms/pkg/utils"
	parse "nms/pkg/utils/jsonParse"
)

var agentsIPs map[byte][4]byte

func StartUDPServer(port string) {
	jsonData := parse.GetDataFromJson("configs/tasks.json")
	var taskList []parse.Task = parse.ParseDataFromJson(jsonData)

	// validate tasks
	parse.ValidateAndRebuildTasks(taskList)

	// Initialize the map
	agentsIPs = make(map[byte][4]byte)

	serverConn := utils.ResolveUDPAddrAndListen("localhost", "8081")
	handleRegistrations(serverConn)

	//serverConn.SetDeadline(time.Now().Add(5 * time.Second))

	//serverConn.Close()

	// go Send tasks

	// go Receive metrics
}
