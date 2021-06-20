package handler

import (
	"encoding/hex"
	"net"
	"regexp"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/miekg/dns"
)

var ipv6regexpsplit = regexp.MustCompile(`(?:\w+-)*?(\d+)-(\d+)-(\d+)-(\d+)\.ipv6\..+$`)
var ipv6regexpHex = regexp.MustCompile(`(?:\w+-)*?([0-9a-f]{32})\.ipv6\..+$`)
var ipv6loopback = net.IPv6loopback

func IPv6Handler(w dns.ResponseWriter, r *dns.Msg) {
	q := r.Question[0]
	m := &dns.Msg{}
	m.SetReply(r)
	defer func() {
		if len(m.Answer) > 0 {
			w.WriteMsg(m)
		}
	}()
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
}
