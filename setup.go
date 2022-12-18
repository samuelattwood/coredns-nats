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

func natsParse(c *caddy.Controller) (*JetStream, error) {
	js := new(JetStream)

	c.Next()
	if c.NextBlock() {
		for {
			switch c.Val() {
			case "url":
				if !c.NextArg() {
					return js, c.ArgErr()
				}
				js.natsURL = c.Val()
			case "bucket":
				if !c.NextArg() {
					return js, c.ArgErr()
				}
				js.bucketName = c.Val()
			default:
				if c.Val() != "}" {
					return js, c.Errf("unknown property %s", c.Val())
				}
			}
			if !c.Next() {
				break
			}
		}
	}

	err := js.Connect()

	return js, err
}
