package main

import (
	"os"
	"fmt"
	"net"
//	"io/ioutil"
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
	conn, err := net.Dial(netType, port)
	checkFatal(err, "Error dialing")
	_, err = conn.Write([]byte("pinggggg"))
	checkFatal(err, "Error writing to connection")
	buf := make([]byte, bufferSize)
	_, err = conn.Read(buf)
//	res, err := ioutil.ReadAll(conn)
	checkFatal(err, "Unable to read response from server")
	fmt.Printf("Response: %s", string(buf))
}