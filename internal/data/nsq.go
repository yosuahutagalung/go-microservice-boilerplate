package data

import (
	"service_boilerplate/internal/conf"

	"github.com/nsqio/go-nsq"
)

func NewNSQProducer(c *conf.Data) (*nsq.Producer, func(), error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(c.Mq.NsqdAddr, config)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() { producer.Stop() }

	return producer, cleanup, nil
}
