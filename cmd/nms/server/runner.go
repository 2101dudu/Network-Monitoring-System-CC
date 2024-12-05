package main

import (
	"flag"
	"log"
	servertcp "nms/internal/server/alertflow"
	serverudp "nms/internal/server/nettask"
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

	// Start tcp server
	go servertcp.StartTCPServer("8080")
	// start udp server
	serverudp.StartUDPServer("8081")
}
