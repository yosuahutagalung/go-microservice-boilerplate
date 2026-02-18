package data

import (
	"context"
	"fmt"
	"service_boilerplate/internal/biz"
	"service_boilerplate/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nsqio/go-nsq"
)

type greeterPublisher struct {
	data          *Data
	producer      *nsq.Producer
	greetingTopic string
	log           *log.Helper
}

func NewGreeterPublisher(c *conf.Data, data *Data, producer *nsq.Producer, logger log.Logger) (biz.GreeterEventPublisher, error) {
	topicInfo, ok := c.Mq.Topics["greeting_events"]
	if !ok {
		return nil, fmt.Errorf("missing NSQ topic configuration for 'greeting_events'")
	}

	return &greeterPublisher{
		data:          data,
		producer:      producer,
		greetingTopic: topicInfo.Topic,
		log:           log.NewHelper(logger),
	}, nil
}

func (p *greeterPublisher) PublishGreetingSaid(ctx context.Context, g *biz.Greeter) error {
	// payload, _ := json.Marshal(g)

	// return p.data.query.InsertOutboxEvent(ctx, db.InsertOutboxEventParams{
	// 	Topic:   p.greetingTopic,
	// 	Payload: payload,
	// })

	return nil
}
