package jsonParse

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetDataFromJson(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("[ERROR 8] Unable to read file: ", err)
		os.Exit(1)
	}

	return data
}

func ParseDataFromJson(data []byte) []Task {
	var tasks []Task
	err := json.Unmarshal(data, &tasks)
	if err != nil {
		fmt.Println("[ERROR 9] Unable to parse data", err)
		os.Exit(1)
	}

	return tasks
}

func ValidateTaskList(taskList []Task) {
	// map to track unique TaskIDs
	seenTaskIDs := make(map[uint16]bool)

	for _, task := range taskList {
		// check if TaskID is repeated
		if seenTaskIDs[task.TaskID] {
			fmt.Printf("[ERROR 20] Duplicate TaskID found: %d\n", task.TaskID)
			os.Exit(1)
		}
		seenTaskIDs[task.TaskID] = true

		valid := validateTask(task)
		if !valid {
			fmt.Println("[ERROR 19] Invalid task")
			os.Exit(1)
		}
	}
}
