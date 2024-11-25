package jsonParse

import (
	"encoding/json"
	"log"
	"os"
)

func GetDataFromJson(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln("[ERROR 8] Unable to read file: ", err)
	}

	return data
}

func ParseDataFromJson(data []byte) []Task {
	var tasks []Task
	err := json.Unmarshal(data, &tasks)
	if err != nil {
		log.Fatalln("[ERROR 9] Unable to parse data:", err)
	}

	return tasks
}

func ValidateTaskList(taskList []Task) {
	// map to track unique TaskIDs
	seenTaskIDs := make(map[uint16]bool)

	for _, task := range taskList {
		// check if TaskID is repeated
		if seenTaskIDs[task.TaskID] {
			log.Fatalf("[ERROR 20] Duplicate TaskID found: %d\n", task.TaskID)
		}
		seenTaskIDs[task.TaskID] = true

		valid := validateTask(task)
		if !valid {
			log.Fatalln("[ERROR 19] Invalid task")
		}
	}
}
