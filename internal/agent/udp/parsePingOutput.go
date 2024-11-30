package udp

func parsePingOutput(output string) string {
	return findInLines("min/", output)
}
