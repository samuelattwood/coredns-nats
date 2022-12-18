package nats

import (
	"github.com/coredns/coredns/plugin"
	"github.com/nats-io/nats.go"
)

type JetStream struct {
	natsURL    string
	bucketName string
	kvBucket   nats.KeyValue
	jetStream  nats.JetStreamContext
	Next       plugin.Handler
}

func (js *JetStream) Connect() error {
	nc, err := nats.Connect(js.natsURL)
	if err != nil {
		return err
	}

	jetStream, err := nc.JetStream()
	if err != nil {
		return err
	}

	kvBucket, err := jetStream.KeyValue(js.bucketName)
	if err != nil {
		return err
	}

	js.jetStream = jetStream
	js.kvBucket = kvBucket

	return nil

}
