package udp

import (
	"log"
	"strconv"
	"strings"
)

func parseIperfOutput(bandwidth bool, jitter bool, packetLoss bool, jitterLimit float64, packetLossLimit float64, output string) (string, bool, bool) {
	jitterHasExceeded := false
	packetLossHasExceeded := false

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

			jitterValue, err := strconv.ParseFloat(separatedLine[8], 64)
			if err != nil {
				log.Println("[AGENT] [ERROR 155] Transforming jitter string into float")
			}
			if jitterValue > jitterLimit { // check if jitter has exceeded
				jitterHasExceeded = true
			}

		}

		if packetLoss {
			newOutput += separatedLine[10] + " " + separatedLine[11]

			packetLossPercentageStr := strings.Trim(separatedLine[11], "()%")
			packetLossPercentage, err := strconv.ParseFloat(packetLossPercentageStr, 64)
			if err != nil {
				log.Println("[AGENT] [ERROR 156] Transforming packet loss percentage string into float")
			}
			if packetLossPercentage > packetLossLimit { // Check if packet loss has exceeded
				packetLossHasExceeded = true
			}
		}

		return newOutput, jitterHasExceeded, packetLossHasExceeded
	}
}
