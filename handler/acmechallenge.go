package handler

import (
	"log"
	"strings"

	"github.com/kawaz/oreore-authority/config"
	"github.com/miekg/dns"
)

func AcmeChallengeHandler(w dns.ResponseWriter, r *dns.Msg) {
	q := r.Question[0]
	m := &dns.Msg{}
	m.SetReply(r)
	if q.Qtype == dns.TypeNS && strings.HasPrefix(q.Name, "_acme-challenge.") {
		for _, d := range config.OreOreConfig.Domains {
			if strings.HasSuffix(q.Name, d.Name) {
				for _, ns := range d.Ns {
					rr := &dns.NS{
						Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 300},
						Ns:  ns,
					}
					m.Answer = append(m.Answer, rr)
				}
			}
		}
		if len(m.Answer) > 0 {
			w.WriteMsg(m)
		}
	}
	if q.Qtype == dns.TypeTXT && strings.HasPrefix(q.Name, "_acme-challenge.") {
		for _, d := range config.OreOreConfig.Domains {
			if strings.HasSuffix(q.Name, d.Name) {
				for _, ns := range d.Ns {
					rr := &dns.NS{
						Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 300},
						Ns:  ns,
					}
					m.Extra = append(m.Extra, rr)
				}
			}
		}
		if len(m.Answer) > 0 {
			w.WriteMsg(m)
		}
	}
	log.Printf("%v, %#v, A:%#v, Ex:%#v", w.RemoteAddr(), r.Question, m.Answer, m.Extra)
}
