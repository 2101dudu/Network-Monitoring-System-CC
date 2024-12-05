package main

import (
	"flag"
	"log"
	agent "nms/internal/agent/nettask"
)

func main() {
	// check if the server was ran with the argument --verbose
	verbose := flag.Bool("verbose", false, "enable verbose mode")
	flag.Parse()
	if *verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	// Start the agent
	agent.StartUDPAgent()

}
