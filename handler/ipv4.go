package handler

import (
	"encoding/hex"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/miekg/dns"
)

var ipv4regexpDigit = regexp.MustCompile(`(?:\w+-)*?(\d+)-(\d+)-(\d+)-(\d+)\.ipv4\..+$`)
var ipv4regexpHex = regexp.MustCompile(`(?:\w+-)*?([0-9a-f]{8})\.ipv4\..+$`)
var ipv4loopback = net.IPv4(127, 0, 0, 1)

func IPv4Handler(w dns.ResponseWriter, r *dns.Msg) {
	q := r.Question[0]
	m := &dns.Msg{}
	m.SetReply(r)
	defer func() {
		if len(m.Answer) > 0 {
			w.WriteMsg(m)
		}
	}()
	if strings.HasPrefix(q.Name, "ipv4.oreore.") && q.Qtype == dns.TypeNS {
		rr := &dns.NS{
			Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 300},
			Ns:  "aws.authority.oreore.net.",
		}
		m.Answer = append(m.Answer, rr)
		return
	}
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
}
