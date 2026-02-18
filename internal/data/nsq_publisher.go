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

type nsqPublisher struct {
	data          *Data
	producer      *nsq.Producer
	greetingTopic string
}

func NewNSQPublisher(c *conf.Data, data *Data, producer *nsq.Producer) (biz.GreeterEventPublisher, error) {
	topicInfo, ok := c.Mq.Topics["greeting_events"]
	if !ok {
		return nil, fmt.Errorf("missing NSQ topic configuration for 'greeting_events'")
	}

	return &nsqPublisher{
		data:          data,
		producer:      producer,
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
