package utils

import (
	"log"
	"net"
)

func WriteTCP(conn *net.TCPConn, data []byte, successMsg string, errorMsg string) {
	// Write the data
	_, err := conn.Write(data)
	if err != nil {
		log.Fatalln(errorMsg, ":", err)
	}

	log.Println(successMsg)
}

func ReadTCP(conn *net.TCPConn, successMessage string, errorMessage string) (int, []byte) {
	newData := make([]byte, BUFFERSIZE)
	n, err := conn.Read(newData)
	if err != nil {
		log.Fatalln(errorMessage, ":", err)
	}
	log.Println(successMessage)
	return n, newData
}

func ResolveTCPAddr(ip string, port string) *net.TCPListener {
	addr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		log.Fatalln("[TCP] [ERROR 1] Unable to resolve address:", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln("[TCP] [ERROR 2] Unable to initialize the server:", err)
	}

	return listener
}

func ResolveTCPAddrAndDial(ip string, port string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		log.Fatalln("[TCP] [ERROR 3] Unable to resolve address:", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln("[TCP] [ERROR 4] Unable to connect:", err)
	}
	return conn
}
