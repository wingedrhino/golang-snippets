package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// RFC 862 says port 7; we do +8000 to move it to userland
const port = ":8007"

// Can also be tcp4 or tcp6
const networkType = "tcp"

// Prevent buffer overruns; unit is bytes
const readLimit = 512

// Read client request in 5 seconds or exit
const readTimeout = 5 * time.Second

func checkFatal(err error, msg string) {
	if err != nil {
		fmt.Printf("Encountered error: %s; %v\n", msg, err)
		os.Exit(0)
	}
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr(networkType, port)
	checkFatal(err, fmt.Sprintf("Unable to ResolveTCPAddr with networkType %s and port %s", networkType, port))

	listener, err := net.ListenTCP(networkType, tcpAddr)
	checkFatal(err, fmt.Sprintf("Unable to ListenTCP at networkType %s and port %s", networkType, port))

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Encountered error while trying to accept connection: %v\n", err)
			continue
		}
		go serve(conn)
	}
}

func serve(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(readTimeout))
	buff := make([]byte, readLimit)
	_, err := conn.Read(buff)

	if err != nil {
		fmt.Printf("Encountered error while reading input: %v\n", err)
		return
	}
	_, err = conn.Write(buff)
	if err != nil {
		fmt.Printf("Encountered error while reading input: %v\n", err)
		return
	}

}
