package main

import (
	"fmt"
	a "nms/src/agent/agent_config"
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
        a.OpenAgent()
    }
}
