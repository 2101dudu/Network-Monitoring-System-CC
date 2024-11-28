package udp

import (
	parse "nms/internal/jsonParse"
)

func handleTasks(taskList []parse.Task) {
	for _, task := range taskList {
		if len(task.Devices) == 1 {
			handlePingTask(task)
		} else {
			handleIperfTask(task)
		}
	}
}
