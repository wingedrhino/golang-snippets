// playing with IP address lookups
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

var ip = flag.String("ip", "", "The IP address you wish to parse. Eg: 192.168.1.1")

func validateFlags() bool {
	flag.Parse()
	if len(*ip) == 0 {
		fmt.Fprintf(os.Stderr, "IP not defined.")
		return false
	}
	return true
}

func main() {
	if validateFlags() == false {
		os.Exit(1)
	}
	addr := net.ParseIP(*ip)
	if addr == nil {
		fmt.Println("Invalid IP")
	} else {
		addrString := addr.String()
		defaultMask := addr.DefaultMask()
		defaultMaskString := defaultMask.String()
		defaultMaskIPString := net.ParseIP(defaultMaskString)
		fmt.Printf("Addess: %s\ndefaultMask: %v\ndefaultMaskString: %s\ndefaultMaskIPString: %s", addrString, defaultMask, defaultMaskString, defaultMaskIPString)
	}
	os.Exit(0)
}
