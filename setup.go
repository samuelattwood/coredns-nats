package nats

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("nats", setup) }

func setup(c *caddy.Controller) error {
	n, err := natsParse(c)
	if err != nil {
		return plugin.Error("nats", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		n.Next = next
		return n
	})

	return nil
}

func natsParse(c *caddy.Controller) (*NATS, error) {
	nats := NATS{}

	c.Next()
	if c.NextBlock() {
		for {
			switch c.Val() {
			case "url":
				if !c.NextArg() {
					return &nats, c.ArgErr()
				}
				nats.URL = c.Val()
			case "bucket":
				if !c.NextArg() {
					return &nats, c.ArgErr()
				}
				nats.Bucket = c.Val()
			default:
				if c.Val() != "}" {
					return &nats, c.Errf("unknown property %s", c.Val())
				}
			}
			if !c.Next() {
				break
			}
		}
	}

	err := nats.Connect()

	return &nats, err
}
