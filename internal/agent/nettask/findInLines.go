package nettask

import "strings"

func findInLines(containsString string, output string) string {
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.Contains(line, containsString) {
			return line
		}
	}

	return ""
}
