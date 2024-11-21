package utils

import (
	"fmt"
	"net"
	"os"
)

func WriteUDP(conn *net.UDPConn, udpAddr *net.UDPAddr, data []byte, successMessage string, errorMessage string) {
	var err error

	if udpAddr == nil { // write UDP without using UDP address
		_, err = conn.Write(data)

	} else { // write UDP using UDP address
		_, err = conn.WriteToUDP(data, udpAddr)
	}

	if err != nil {
		fmt.Println(errorMessage, ":", err)
		os.Exit(1)
	}
	fmt.Println(successMessage)
}

func ReadUDP(conn *net.UDPConn, successMessage string, errorMessage string) (int, *net.UDPAddr, []byte) {
	newData := make([]byte, 1024)
	n, udpAddr, err := conn.ReadFromUDP(newData)
	if err != nil {
		fmt.Println(errorMessage, ":", err)
		os.Exit(1)
	}
	fmt.Println(successMessage)
	return n, udpAddr, newData
}

func ResolveUDPAddrAndListen() *net.UDPConn {
	ip := getIPAddress()
	addr, err := net.ResolveUDPAddr("udp", ip+":9091")
	if err != nil {
		fmt.Println("[AGENT] [ERROR 8] Unable to resolve address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("[AGENT] [ERROR 9] Unable to initialize the agent:", err)
		os.Exit(1)
	}

	return conn
}
