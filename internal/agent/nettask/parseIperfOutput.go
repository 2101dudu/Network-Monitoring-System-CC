package nettask

import (
	"log"
	utils "nms/internal/utils"
	"strconv"
	"strings"
)

func parseIperfOutput(bandwidth bool, jitter bool, packetLoss bool, jitterLimit float32, packetLossLimit float32, output string) (string, float32, float32) {
	jitterHasExceeded := float32(0)
	packetLossHasExceeded := float32(0)
	shift := 0

	if len(output) == 0 {
		return "", jitterHasExceeded, packetLossHasExceeded
	}

	if bandwidth {
		line := findInLines("sec", output)
		if len(line) == 0 {
			return "", jitterHasExceeded, packetLossHasExceeded
		}

		if len(line) == 8 {
			shift = 1
		}
		separatedLine := strings.Fields(line)
		return separatedLine[6+shift] + " " + separatedLine[7+shift], jitterHasExceeded, packetLossHasExceeded
	} else {
		line := findInLines("%", output)
		if len(line) == 0 {
			return "", jitterHasExceeded, packetLossHasExceeded
		}

		if len(line) == 11 || len(line) == 14 {
			shift = 1
		}

		separatedLine := strings.Fields(line)

		newOutput := ""
		if jitter {
			newOutput += separatedLine[8+shift] + " " + separatedLine[9+shift] + " "

			jitterValue, err := strconv.ParseFloat(separatedLine[8+shift], 32)
			if err != nil {
				log.Println(utils.Red+"[ERROR 165] Transforming jitter string into float", utils.Reset)
			}
			if float32(jitterValue) > jitterLimit { // check if jitter has exceeded
				jitterHasExceeded = float32(jitterValue) - float32(jitterLimit)
			}

		}

		if packetLoss {
			newOutput += separatedLine[10+shift] + separatedLine[11+shift] + " " + separatedLine[12+shift]

			packetLossPercentageStr := strings.Trim(separatedLine[12+shift], "()%")
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
