package handler

import (
	"regexp"

	"github.com/miekg/dns"
)

var localHostnameRegex = regexp.MustCompile(`^([\w\-]+\.local\.).+$`)

func LocalHandler(w dns.ResponseWriter, r *dns.Msg) {
	// (hostname.local.)XXXX 部分を取り出す
	q := r.Question[0]
	sub := localHostnameRegex.FindSubmatch([]byte(q.Name))
	if sub == nil {
		return
	}
	hostname := sub[1]

	m := new(dns.Msg)
	m.SetReply(r)
	rr := &dns.CNAME{
		Hdr:    dns.RR_Header{Name: q.Name, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 24 * 3600},
		Target: dns.Fqdn(string(hostname)),
	}
	m.Answer = append(m.Answer, rr)
	if len(m.Answer) > 0 {
		w.WriteMsg(m)
	}
}
