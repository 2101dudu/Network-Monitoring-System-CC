package main

import (
	sc "nms/internal/server/server_config"
	u "nms/pkg/utils"
	"path/filepath"
)

func main() {
	// get path to json file
	filePath := filepath.Join("configs", "settings.json")

	// get []byte structure corresponding to json data
	jsonData := u.GetDataFromJson(filePath)

	// get []Task from json data
	tasksList := u.ParseDataFromJson(jsonData)

	sc.OpenServer(tasksList)
}
