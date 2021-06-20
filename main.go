package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kawaz/oreore-authority/handler"
	"github.com/miekg/dns"
)

type DomainConfig struct {
	Name string
	Ns   []string
}

var DomainConfigs []*DomainConfig = []*DomainConfig{
	{
		Name: "oreore.net.",
		Ns: []string{
			"ns-973.awsdns-57.net.",
			"ns-409.awsdns-51.com.",
			"ns-1025.awsdns-00.org.",
			"ns-1854.awsdns-39.co.uk.",
		},
	},
	{
		Name: "oreore.dev.",
		Ns:   []string{},
	},
	{
		Name: "oreore.app.",
		Ns:   []string{},
	},
	{
		Name: "oreore.page.",
		Ns:   []string{},
	},
}

func main() {
	for _, d := range DomainConfigs {
		dns.HandleFunc("ipv4."+d.Name, handler.IPv4Handler)
		dns.HandleFunc("ipv6."+d.Name, handler.IPv6Handler)
		dns.HandleFunc("local."+d.Name, handler.LocalHandler)
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
