package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
)

const networkType = "tcp"

var address = flag.String("address", "http://www.google.com", "The address/URL you wish to fetch. Should be a valid HTTP address. Eg: http://en.wikipedia.org/wiki/Van_Halen")

// checkErr checks if error is nil, if not prints message and quits.
func checkErr(err error, msg string) {
	if err != nil {
		fmt.Printf("Encountered Error: %s\n%v\n", msg, err)
		os.Exit(1)
	}
}

func initFlags() {
	flag.Parse()

	if len(*address) == 0 {
		fmt.Printf("Address cannot be blank!\n")
		os.Exit(1)
	}
}

func main() {
	initFlags()

	u, err := url.Parse(*address)
	checkErr(err, "Unable to parse URL")

	var serverPort = u.Host
	// If no port is passed, assume port 80
	if len(u.Port()) == 0 {
		serverPort += ":80"
	}

	tcpAddr, err := net.ResolveTCPAddr(networkType, serverPort)
	checkErr(err, "Unable to resolve address")

	conn, err := net.DialTCP(networkType, nil, tcpAddr)
	checkErr(err, "Unable to connect to server")

	escapedPath := u.EscapedPath()
	if len(escapedPath) == 0 {
		escapedPath = "/"
	}
	// Note that we add a Host header here so some servers that use virtual
	// hosting parse the request correctly.
	request := fmt.Sprintf("GET %s HTTP/1.0\r\n\r\nHost: %s", escapedPath, u.Host)
	_, err = conn.Write([]byte(request))
	checkErr(err, "Unable to write response to server.")

	responseBytes, err := ioutil.ReadAll(conn)
	checkErr(err, "Unable to read response from server.")
	response := string(responseBytes)

	fmt.Printf("Response from server:\n%s\n", response)
	os.Exit(0)
}
