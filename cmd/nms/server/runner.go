package main

import (
	"log"
	servertcp "nms/internal/server/tcp"
	serverudp "nms/internal/server/udp"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})

	// Start tcp server
	go func() {
		//log.Println("[MAIN] Starting TCP Server on port 8080...")
		servertcp.StartTCPServer("8080", done)
	}()

	// start udp server
	//log.Println("[MAIN] Starting UDP Server on port 8081...")
	serverudp.StartUDPServer("8081")

	// wait for interrupting signal
	<-quit
	log.Println("[MAIN] Shutting down servers...")

	// Close tcp server
	close(done)

	log.Println("[MAIN] Shutdown complete.")
}
