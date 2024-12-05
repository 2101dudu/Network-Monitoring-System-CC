package nettask

import (
	"bufio"
	"fmt"
	"log"
	"nms/internal/utils"
	"os/exec"
	"strconv"
	"strings"
)

func askNumAgents(reader *bufio.Reader) int {
	for {
		// ask the user for the number of agents
		fmt.Print(utils.Blue+"Enter the number of agents: ", utils.Reset)
		input := parseString(reader)

		// convert the input to an integer
		numAgents, err := strconv.Atoi(input)

		if err == nil && numAgents >= 1 && numAgents <= 255 {
			return numAgents
		}
	}
}

func askJsonPath(reader *bufio.Reader) string {
	// ask the user for the path to the JSON file
	fmt.Print(utils.Blue+"Enter the path to the tasks file: ", utils.Reset)
	return parseString(reader)
}

func parseString(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(utils.Red+"[ERROR 1] Unable to read input", utils.Reset)
	}

	// trimm the input
	return strings.TrimSpace(input)
}

func consultMetricsFile() {
	cmd := exec.Command("cat", "output/metrics.json")

	// run the command
	output, err := cmd.CombinedOutput()

	fmt.Print(string(output))

	if err != nil {
		log.Println(utils.Red+"[ERROR 21] Unable to consult metrics file", utils.Reset)
	}
}

func consultAlertsFile() {
	cmd := exec.Command("cat", "output/alerts.json")

	// run the command
	output, err := cmd.CombinedOutput()

	fmt.Print(string(output))

	if err != nil {
		log.Println(utils.Red+"[ERROR 22] Unable to consult alerts file", utils.Reset)
	}
}
