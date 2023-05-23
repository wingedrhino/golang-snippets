package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const port = ":8013"
const tcpType = "tcp"

// checkErr checks if error is nil, if not prints message and quits.
func checkErr(err error, msg string) {
	if err != nil {
		fmt.Printf("Encountered Error: %s\n%v\n", msg, err)
		os.Exit(1)
	}
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr(tcpType, port)
	checkErr(err, "Unable to resolve TCP address")

	listener, err := net.ListenTCP(tcpType, tcpAddr)
	checkErr(err, "Unable to listen at specified address")

	// single threaded server ahead
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Encountered error while listening: %v\n", err)
			continue
		}

		daytime := time.Now().String()

		fmt.Printf("Writing daytime: %s\n", daytime)
		_, err = conn.Write([]byte(daytime))
		if err != nil {
			fmt.Printf("Error writing daytime: %s\n", daytime)
		}

		err = conn.Close()
		if err != nil {
			fmt.Printf("A connection at daytime %s was unable to be closed.\n", err)
		}
	}
}
