// playing with DNS lookups
package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/miekg/dns"
)

var domain = flag.String("domain", "", "The domain name you wish to lookup. Eg: google.com")

func validateFlags() bool {
	flag.Parse()
	if len(*domain) == 0 {
		fmt.Fprintf(os.Stderr, "IP not defined.")
		return false
	}
	return true
}

func main() {
	if validateFlags() == false {
		os.Exit(1)
	}
	addrRegular, err := net.LookupHost(*domain)
	if err != nil {
		fmt.Printf("Net Package DNS Lookup failed with error: %v\n", err)
	} else {
		fmt.Printf("Net Package DNS Search Result: %v\n,result count: %d\n", addrRegular, len(addrRegular))
	}
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(*domain), dns.TypeA)
	c := dns.Client{}
	addrCustom, rtt, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		fmt.Printf("Custom DNS Lookup failed with error: %v\n", err)
	} else {
		fmt.Printf("Custom DNS Search Result: %v\n,rtt: %d\n", addrCustom, rtt)
	}
	os.Exit(0)
}
