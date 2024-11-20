package utils

import (
	"bytes"
	"os/exec"
	"runtime"
)

// Executes a command and returns the output (and error if exists)
func ExecuteCommand(command string) (string, error) {
	var cmd *exec.Cmd

	// Detect OS (not really needed I guess - we will be working with core linux - but its an idea that ocurred)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command) // Windows
	} else {
		cmd = exec.Command("sh", "-c", command) // Linux/Mac
	}

	// buffer to store the output from standart output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output //we can other buffer for the standart error maybe

	// run
	err := cmd.Run()

	return output.String(), err
}
