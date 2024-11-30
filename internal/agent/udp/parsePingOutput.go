package udp

func parsePingOutput(output string) string {
	return findInLines([]string{"min/"}, output)
}
