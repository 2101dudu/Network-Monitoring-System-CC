package udp

import (
	"log"
	"strconv"
	"strings"
)

func parseIperfOutput(bandwidth bool, jitter bool, packetLoss bool, jitterLimit float32, packetLossLimit float32, output string) (string, float32, float32) {
	jitterHasExceeded := float32(0)
	packetLossHasExceeded := float32(0)

	if bandwidth {
		line := findInLines("sec", output)
		separatedLine := strings.Fields(line)
		return separatedLine[6] + " " + separatedLine[7], jitterHasExceeded, packetLossHasExceeded
	} else {
		line := findInLines("%", output)
		separatedLine := strings.Fields(line)

		newOutput := ""
		if jitter {
			newOutput += separatedLine[8] + " " + separatedLine[9] + " "

			jitterValue, err := strconv.ParseFloat(separatedLine[8], 32)
			if err != nil {
				log.Println("[ERROR 155] Transforming jitter string into float")
			}
			if float32(jitterValue) > jitterLimit { // check if jitter has exceeded
				jitterHasExceeded = float32(jitterValue) - float32(jitterLimit)
			}

		}

		if packetLoss {
			newOutput += separatedLine[10] + " " + separatedLine[11]

			packetLossPercentageStr := strings.Trim(separatedLine[11], "()%")
			packetLossPercentage, err := strconv.ParseFloat(packetLossPercentageStr, 32)
			if err != nil {
				log.Println("[ERROR 156] Transforming packet loss percentage string into float")
			}
			if float32(packetLossPercentage) > packetLossLimit { // Check if packet loss has exceeded
				packetLossHasExceeded = float32(packetLossPercentage) - float32(packetLossLimit)
			}
		}

		return newOutput, jitterHasExceeded, packetLossHasExceeded
	}
}
