package nettask

import (
	"fmt"
	"log"
	utils "nms/internal/utils"
	"strconv"
	"strings"
)

func parseIperfOutput(bandwidth bool, jitter bool, packetLoss bool, jitterLimit float32, packetLossLimit float32, output string) (string, float32, float32) {
	jitterHasExceeded := float32(0)
	packetLossHasExceeded := float32(0)
	fmt.Println(output)

	if len(output) == 0 {
		return "", jitterHasExceeded, packetLossHasExceeded
	}

	if bandwidth {
		line := findInLines("sec", output)
		separatedLine := strings.Fields(line)
		return separatedLine[7] + " " + separatedLine[8], jitterHasExceeded, packetLossHasExceeded
	} else {
		line := findInLines("%", output)
		separatedLine := strings.Fields(line)

		newOutput := ""
		if jitter {
			newOutput += separatedLine[9] + " " + separatedLine[10] + " "

			jitterValue, err := strconv.ParseFloat(separatedLine[9], 32)
			if err != nil {
				log.Println(utils.Red+"[ERROR 165] Transforming jitter string into float", utils.Reset)
			}
			if float32(jitterValue) > jitterLimit { // check if jitter has exceeded
				jitterHasExceeded = float32(jitterValue) - float32(jitterLimit)
			}

		}

		if packetLoss {
			newOutput += separatedLine[11] + separatedLine[12] + " " + separatedLine[13]

			packetLossPercentageStr := strings.Trim(separatedLine[13], "()%")
			packetLossPercentage, err := strconv.ParseFloat(packetLossPercentageStr, 32)
			if err != nil {
				log.Println(utils.Red+"[ERROR 166] Transforming packet loss percentage string into float", utils.Reset)
			}
			if float32(packetLossPercentage) > packetLossLimit { // Check if packet loss has exceeded
				packetLossHasExceeded = float32(packetLossPercentage) - float32(packetLossLimit)
			}
		}

		return newOutput, jitterHasExceeded, packetLossHasExceeded
	}
}
