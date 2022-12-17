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

func (handler *NATS) Name() string { return name }

func (handler *NATS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true
	a.Compress = true

	kv := *handler.KV

	domain := strings.TrimSuffix(state.QName(), ".")

	jsKey, err := kv.Get(domain)
	if err != nil {
		return dns.RcodeServerFailure, nil
	}

	var jsRecord A_Record
	json.Unmarshal(jsKey.Value(), &jsRecord)

	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{
		Name:   dns.Fqdn(domain),
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    360,
	}
	rr.A = jsRecord.Ip

	a.Answer = append(a.Answer, rr)

	state.Scrub(a)
	w.WriteMsg(a)

	return dns.RcodeSuccess, nil
}
