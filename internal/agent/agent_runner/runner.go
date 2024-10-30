package main

import (
	"fmt"
	ac "nms/internal/agent/agent_config"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: ./agent.go -<option>")
		os.Exit(1)
	}

	switch args[0] {
	case "-udp":
	case "-tcp":
		ac.OpenAgent()
	}
}
