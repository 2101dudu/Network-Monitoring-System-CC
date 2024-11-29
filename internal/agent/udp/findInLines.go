package udp

import "strings"

func findInLines(containsStrings []string, output string) string {
	lines := strings.Split(output, "\n")

	final := ""

	for _, line := range lines {
		for _, containsString := range containsStrings {
			if strings.Contains(line, containsString) {
				final += line + "\n"
			}
		}
	}

	return final
}
