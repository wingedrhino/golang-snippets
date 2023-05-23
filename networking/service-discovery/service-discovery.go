package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

var networkType = flag.String("networkType", "tcp", "The network type. Can be tcp or udp. Defaults to tcp.")
var serviceName = flag.String("serviceName", "", "Name of the service you wish to look up the port for. Eg: dns")

// validateFlags checks and validates set flags
func validateFlags() bool {
	flag.Parse()
	if len(*serviceName) == 0 {
		fmt.Printf("Expected flag service to be set.\n")
		return false
	}
	return true
}

func main() {
	if validateFlags() != true {
		os.Exit(1)
	}

	port, err := net.LookupPort(*networkType, *serviceName)

	if err != nil {
		fmt.Printf("Lookup failed with error: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("Port: %d\n", port)
}
