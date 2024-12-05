package utils

import (
	"log"
	"net"
)

func WriteTCP(conn *net.TCPConn, data []byte, alertMsg string, errorMsg string) {
	// Write the data
	_, err := conn.Write(data)
	if err != nil {
		log.Fatalln(Red+errorMsg, ":", err, Reset)
	}

	log.Println(Magenta+alertMsg, Reset)
}

func ReadTCP(conn *net.TCPConn, alertMsg string, errorMessage string) (int, []byte) {
	newData := make([]byte, BUFFERSIZE)
	n, err := conn.Read(newData)
	if err != nil {
		log.Fatalln(Red+errorMessage, ":", err, Reset)
	}
	log.Println(Magenta+alertMsg, Reset)
	return n, newData
}

func ResolveTCPAddr(ip string, port string) *net.TCPListener {
	addr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		log.Fatalln(Red+"[ERROR 1] Unable to resolve address:", err, Reset)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(Red+"[ERROR 2] Unable to initialize the server:", err, Reset)
	}

	return listener
}

func ResolveTCPAddrAndDial(ip string, port string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		log.Fatalln(Red+"[ERROR 3] Unable to resolve address:", err, Reset)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln(Red+"[ERROR 4] Unable to connect:", err, Reset)
	}
	return conn
}
