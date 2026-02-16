package data

import (
	"context"
	"encoding/json"
	"fmt"
	"service_boilerplate/internal/biz"
	"service_boilerplate/internal/conf"
	"service_boilerplate/internal/data/db"

	"github.com/nsqio/go-nsq"
)

// 1. Extract the raw Producer connection into its own Wire Provider
func NewNSQProducer(c *conf.Data) (*nsq.Producer, func(), error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(c.Mq.NsqdAddr, config)
	if err != nil {
		return nil, nil, err
	}

	// This cleanup function will close the shared connection when Kratos shuts down
	cleanup := func() { producer.Stop() }

	return producer, cleanup, nil
}

type nsqPublisher struct {
	data          *Data
	producer      *nsq.Producer // We still keep the reference
	greetingTopic string
}

// 2. Modify this to ACCEPT the *nsq.Producer as a parameter
// Notice we also removed the `func()` cleanup return, because NewNSQProducer handles that now!
func NewNSQPublisher(c *conf.Data, data *Data, producer *nsq.Producer) (biz.GreeterEventPublisher, error) {
	topicInfo, ok := c.Mq.Topics["greeting_events"]
	if !ok {
		return nil, fmt.Errorf("missing NSQ topic configuration for 'greeting_events'")
	}

	return &nsqPublisher{
		data:          data,
		producer:      producer, // Use the injected producer
		greetingTopic: topicInfo.Topic,
	}, nil
}

func (p *nsqPublisher) PublishGreetingSaid(ctx context.Context, g *biz.Greeter) error {
	payload, _ := json.Marshal(g)

	return p.data.query.InsertOutboxEvent(ctx, db.InsertOutboxEventParams{
		Topic:   p.greetingTopic,
		Payload: payload,
	})
}
