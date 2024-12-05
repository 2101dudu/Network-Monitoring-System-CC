package nettask

func parsePingOutput(output string) string {
	return findInLines("min/", output)
}
