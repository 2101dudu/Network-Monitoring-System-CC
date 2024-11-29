package udp

import "strings"

func parseIperfOutput(bandwidth bool, jitter bool, packetLoss bool, output string) string {
	if bandwidth {
		line := findInLines([]string{"sec"}, output)
		separatedLine := strings.Fields(line)
		return separatedLine[6] + " " + separatedLine[7]
	} else {
		line := findInLines([]string{"%"}, output)
		separatedLine := strings.Fields(line)

		newOutput := ""
		if jitter {
			newOutput += separatedLine[8] + " " + separatedLine[9] + " "
		}
		if packetLoss {
			newOutput += separatedLine[10] + " " + separatedLine[11]
		}

		return newOutput
	}
}
