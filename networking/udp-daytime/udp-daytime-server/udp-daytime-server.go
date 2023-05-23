package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Port we're alistening at
const port = ":8013"

// May be "udp", "udp4" or "udp6"
const netType = "udp"

// Size in bytes
const bufferSize = 512

// checkFatal checks if error is nil, if not prints message and quits.
func checkFatal(err error, msg string) {
	if err != nil {
		fmt.Printf("Encountered Error: %s\n%v\n", msg, err)
		os.Exit(1)
	}
}

func main() {
	udpAddr, err := net.ResolveUDPAddr(netType, port)
	checkFatal(err, "Unable to ResolveUDPAddr")

	conn, err := net.ListenUDP(netType, udpAddr)
	defer conn.Close()
	checkFatal(err, "Unable to ListenUDP")

	readBuffer := make([]byte, bufferSize)

	for {
		inLen, clAddr, err := conn.ReadFromUDP(readBuffer)
		fmt.Printf("Incoming Message from client %v: %s\n", *clAddr, string(readBuffer))
		go serve(conn, inLen, clAddr, err)
	}
}

func serve(conn *net.UDPConn, inLen int, clAddr *net.UDPAddr, err error) {
	if inLen == 0 {
		fmt.Printf("0 input length from client %v\n", *clAddr)
		return
	}
	if err != nil {
		fmt.Printf("Error encountered while reading from connection: %v\n", err)
		return
	}

	// conn, err := net.DialUDP(netType, nil, clAddr)
	// if err != nil {
	// 	fmt.Printf("Error opening connection to client %v\n", *clAddr)
	// 	return
	// }
	t := time.Now().String()
	fmt.Printf("Responding to %v with time %s\n", *clAddr, t)
	res := []byte(t)
	_, err = conn.WriteToUDP(res, clAddr)
	if err != nil {
		fmt.Printf("Error responding to %v with time %s: %v\n", *clAddr, t, err)
	}
}
