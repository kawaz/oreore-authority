package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	"github.com/k0kubun/pp"
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
		dns.HandleFunc("ipv4."+fqdn, handlerIPv4)
		dns.HandleFunc("ipv6."+fqdn, handlerIPv6)
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

var ipv4regexpDigit = regexp.MustCompile(`(?:\w+-)*?(\d+)-(\d+)-(\d+)-(\d+)\.ipv4\..+$`)
var ipv4regexpHex = regexp.MustCompile(`(?:\w+-)*?([0-9a-f]{8})\.ipv4\..+$`)
var ipv4loopback = net.IPv4(127, 0, 0, 1)

func handlerIPv4(w dns.ResponseWriter, r *dns.Msg) {
	q := r.Question[0]
	var ip net.IP
	sub := ipv4regexpDigit.FindSubmatch([]byte(q.Name))
	if sub != nil {
		ip1, _ := strconv.Atoi(string(sub[1]))
		ip2, _ := strconv.Atoi(string(sub[2]))
		ip3, _ := strconv.Atoi(string(sub[3]))
		ip4, _ := strconv.Atoi(string(sub[4]))
		if ip1 < 256 && ip2 < 256 && ip3 < 256 && ip4 < 256 {
			ip = net.IPv4(byte(ip1), byte(ip2), byte(ip3), byte(ip4))
		}
	}
	if ip == nil {
		sub := ipv4regexpHex.FindSubmatch([]byte(q.Name))
		if sub != nil {
			pp.Println(sub, ip)
			bytes, err := hex.DecodeString(string(sub[1]))
			if len(bytes) == net.IPv4len && err == nil {
				ip = net.IP(bytes)
			}
		}
	}
	if ip == nil {
		ip = ipv4loopback
	}

	m := new(dns.Msg)
	m.SetReply(r)
	switch q.Qtype {
	case dns.TypeA:
		rr := &dns.A{
			Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 24 * 3600},
			A:   ip.To4(),
		}
		m.Answer = append(m.Answer, rr)
	case dns.TypeAAAA:
		rr := &dns.AAAA{
			Hdr:  dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 24 * 3600},
			AAAA: ip.To16(),
		}
		m.Answer = append(m.Answer, rr)
	}
	if len(m.Answer) > 0 {
		w.WriteMsg(m)
	}
}

var ipv6regexpsplit = regexp.MustCompile(`(?:\w+-)*?(\d+)-(\d+)-(\d+)-(\d+)\.ipv6\..+$`)
var ipv6regexpHex = regexp.MustCompile(`(?:\w+-)*?([0-9a-f]{32})\.ipv6\..+$`)
var ipv6loopback = net.IPv6loopback

func handlerIPv6(w dns.ResponseWriter, r *dns.Msg) {
	q := r.Question[0]
	var ip net.IP
	if ip == nil {
		sub := ipv6regexpHex.FindSubmatch([]byte(q.Name))
		if sub != nil {
			pp.Println(sub, ip)
			bytes, err := hex.DecodeString(string(sub[1]))
			if len(bytes) == net.IPv6len && err == nil {
				ip = net.IP(bytes)
			}
		}
	}
	if ip == nil {
		ip = net.IPv6loopback
	}

	m := new(dns.Msg)
	m.SetReply(r)
	switch q.Qtype {
	case dns.TypeA:
		ipv4 := ip.To4()
		if ipv4 == nil {
			ipv4 = ipv4loopback
		}
		rr := &dns.A{
			Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 24 * 3600},
			A:   ipv4,
		}
		m.Answer = append(m.Answer, rr)
	case dns.TypeAAAA:
		rr := &dns.AAAA{
			Hdr:  dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 24 * 3600},
			AAAA: ip,
		}
		m.Answer = append(m.Answer, rr)
	}
	if len(m.Answer) > 0 {
		w.WriteMsg(m)
	}
}
