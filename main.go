package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kawaz/oreore-resolver/handler"
	"github.com/miekg/dns"
)

func main() {
	fqdns := []string{
		"oreore.net.",
		"oreore.dev.",
		"oreore.app.",
		"oreore.page.",
	}
	for _, fqdn := range fqdns {
		dns.HandleFunc("ipv4."+fqdn, handler.IPv4Handler)
		dns.HandleFunc("ipv6."+fqdn, handler.IPv6Handler)
		dns.HandleFunc("local."+fqdn, handler.LocalHandler)
	}
	var name, secret string
	go serve("tcp", name, secret, false)
	go serve("udp", name, secret, false)
	// シグナル来るまで止めない
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping\n", s)
}

func serve(net, name, secret string, soreuseport bool) {
	switch name {
	case "":
		server := &dns.Server{Addr: "[::]:53", Net: net, TsigSecret: nil, ReusePort: soreuseport}
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to setup the "+net+" server: %s\n", err.Error())

		}
	default:
		server := &dns.Server{Addr: ":53", Net: net, TsigSecret: map[string]string{name: secret}, ReusePort: soreuseport}
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to setup the "+net+" server: %s\n", err.Error())
		}
	}
}
