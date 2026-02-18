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
	data     *Data
	producer *nsq.Producer
	topics   map[string]*conf.Data_MQ_TopicInfo
	log      *log.Helper
}

func NewGreeterPublisher(c *conf.Data, data *Data, producer *nsq.Producer, logger log.Logger) (biz.GreeterEventPublisher, error) {
	return &greeterPublisher{
		data:     data,
		producer: producer,
		topics:   c.Mq.Topics,
		log:      log.NewHelper(logger),
	}, nil
}

func (p *greeterPublisher) PublishGreetingSaid(ctx context.Context, g *biz.Greeter) error {
	_, exists := p.topics["greeting_events"]
	if !exists {
		return fmt.Errorf("missing NSQ topic configuration for 'greeting_events'")
	}

	// payload, _ := json.Marshal(g)

	// return p.data.query.InsertOutboxEvent(ctx, db.InsertOutboxEventParams{
	// 	Topic:   p.greetingTopic,
	// 	Payload: payload,
	// })

	return nil
}
