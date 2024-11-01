package utils

import (
	"fmt"
	"net"
	"os"
)

func WriteUDP(conn *net.UDPConn, data []byte, successMessage string, errorMessage string) {
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println(errorMessage, ":", err)
		os.Exit(1)
	}
	fmt.Println(successMessage)
}

func ReadUDP(conn *net.UDPConn, successMessage string, errorMessage string) (int, []byte) {
	newData := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(newData)
	if err != nil {
		fmt.Println(errorMessage, ":", err)
		os.Exit(1)
	}
	fmt.Println(successMessage)
	return n, newData
}
