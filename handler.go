package nats

import (
	"encoding/json"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"

	"golang.org/x/net/context"
)

const (
	name = "nats"
)

func (js *JetStream) Name() string { return name }

func (js *JetStream) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = true

	domain := reverseDomain(state.QName())

	kvEntry, err := js.kvBucket.Get(domain)
	if err != nil {
		return dns.RcodeServerFailure, nil
	}

	rr := getType(state.QType())

	err = json.Unmarshal(kvEntry.Value(), &rr)
	if err != nil {
		return dns.RcodeServerFailure, nil
	}

	m.Answer = append(m.Answer, rr)

	state.Scrub(m)
	w.WriteMsg(m)

	return dns.RcodeSuccess, nil
}

func reverseDomain(domain string) string {
	segments := strings.Split(domain, ".")
	for i, j := 0, len(segments)-1; i < j; i, j = i+1, j-1 {
		segments[i], segments[j] = segments[j], segments[i]
	}
	rev := strings.Join(segments, ".")
	return strings.TrimPrefix(rev, ".")
}

func getType(qType uint16) dns.RR {
	switch qType {
	case dns.TypeA:
		return new(dns.A)
	case dns.TypeAAAA:
		return new(dns.AAAA)
	case dns.TypeCNAME:
		return new(dns.CNAME)
	case dns.TypeMX:
		return new(dns.MX)
	case dns.TypeNS:
		return new(dns.NS)
	case dns.TypePTR:
		return new(dns.PTR)
	case dns.TypeSOA:
		return new(dns.SOA)
	case dns.TypeTXT:
		return new(dns.TXT)
	default:
		// Return nil if the specified record type is not supported
		return new(dns.A)
	}

}
