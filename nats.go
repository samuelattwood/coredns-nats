package nats

import (
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"github.com/nats-io/nats.go"
)

type NATS struct {
	URL    string
	Bucket string
	KV     *nats.KeyValue
	Next   plugin.Handler
}

func (handler *NATS) A(name string, record A_Record) (answers, extras []dns.RR) {
	if record.Ip == nil {
		return
	}

	r := new(dns.A)
	r.Hdr = dns.RR_Header{
		Name:   dns.Fqdn(name),
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    360,
	}
	r.A = record.Ip
	answers = append(answers, r)

	return

}

func (handler *NATS) Connect() error {
	nc, err := nats.Connect(handler.URL)
	if err != nil {
		return err
	}

	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	kv, err := js.KeyValue(handler.Bucket)
	if err != nil {
		return err
	}

	handler.KV = &kv

	return nil

}
