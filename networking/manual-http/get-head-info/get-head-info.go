package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

var serviceAddr = flag.String("serviceAddr", "www.example.com:80", "Address of the service. Eg: 'www.example.com:80'")
var networkType = flag.String("networkType", "tcp", "'tcp4', 'tcp6' or 'tcp' (use any). Defaults to 'tcp'")
var lineTerminator = flag.String("lineTerminator", "windows", "'windows' (CRLF; default) or 'unix' (LF).")

func initFlags() {
	flag.Parse()

	if len(*serviceAddr) == 0 {
		fmt.Printf("You haven't entered a URL!")
		os.Exit(1)
	}

	if *networkType != "tcp" && *networkType != "tcp4" && *networkType != "tcp6" {
		fmt.Printf("networkType should be tcp4, tcp6 or tcp")
		os.Exit(1)
	}

	// Set lineTerminator to the actual value based on what's in the flag
	switch *lineTerminator {
	case "windows":
		*lineTerminator = "\r\n"
	case "unix":
		*lineTerminator = "\n"
	default:
		fmt.Printf("Line Terminator should be either 'windows' or 'unix'")
		os.Exit(1)
	}
}

// checkErr checks if error is nil, if not prints message and quits.
func checkErr(err error, msg string) {
	if err != nil {
		fmt.Printf("Encountered Error: %s\n%v\n", msg, err)
		os.Exit(1)
	}
}

func main() {
	initFlags()

	tcpAddr, err := net.ResolveTCPAddr(*networkType, *serviceAddr)
	checkErr(err, "Unable to parse URL")

	conn, err := net.DialTCP(*networkType, nil, tcpAddr)
	checkErr(err, "Unable to connect to server")

	request := "HEAD / HTTP/1.0" + *lineTerminator + *lineTerminator
	_, err = conn.Write([]byte(request))
	checkErr(err, "Unable to write message to server")

	responseBytes, err := ioutil.ReadAll(conn)
	checkErr(err, "Unable to read response from server")
	response := string(responseBytes)

	fmt.Printf("Response from server:\n%s\n", response)

	os.Exit(0)
}
