package jsonParse

import (
	"encoding/json"
	"log"
	"nms/internal/utils"
	"os"
)

func GetDataFromJson(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln(utils.Red+"[ERROR 8] Unable to read file: ", err, utils.Reset)
	}

	return data
}

func ParseDataFromJson(data []byte) []Task {
	log.Println(utils.Blue+"Parsing data from tasks file...", utils.Reset)

	var tasks []Task
	err := json.Unmarshal(data, &tasks)
	if err != nil {
		log.Fatalln(utils.Red+"[ERROR 9] Unable to parse data:", err, utils.Reset)
	}

	log.Println(utils.Green+"Data parsed successfully!", utils.Reset)

	return tasks
}

func ValidateTaskList(taskList []Task) {
	log.Println(utils.Blue+"Validating tasks...", utils.Reset)

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
			log.Fatalln(utils.Red+"[ERROR 19] Invalid task", utils.Reset)
		}
	}

	log.Println(utils.Green+"Tasks validated successfully!", utils.Reset)
}
