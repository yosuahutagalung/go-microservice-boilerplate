package server

import (
	"context"

	"service_boilerplate/internal/conf"
	"service_boilerplate/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nsqio/go-nsq"
)

// NSQServer manages all background message consumers
type NSQServer struct {
	consumers []*nsq.Consumer
	lookupd   string
	log       *log.Helper
}

func NewNSQServer(c *conf.Data, svc *service.GreeterService, logger log.Logger) (*NSQServer, error) {
	srv := &NSQServer{
		lookupd: c.Mq.LookupdAddr,
		log:     log.NewHelper(logger),
	}

	nsqConfig := nsq.NewConfig()

	if info, exists := c.Mq.Topics["greeting_events"]; exists {
		consumer, err := nsq.NewConsumer(info.Topic, info.Channel, nsqConfig)
		if err != nil {
			return nil, err
		}

		// Attach the handler from the service
		consumer.AddHandler(nsq.HandlerFunc(svc.HandleGreetingEvent))

		// Add to our list of managed consumers
		srv.consumers = append(srv.consumers, consumer)
	} else {
		srv.log.Warn("⚠️ WARNING: 'greeting_events' missing from config.yaml topics map")
	}

	return srv, nil
}

// Start opens the network connections to NSQ when Kratos boots up
func (s *NSQServer) Start(ctx context.Context) error {
	for _, consumer := range s.consumers {
		err := consumer.ConnectToNSQLookupd(s.lookupd)
		if err != nil {
			return err
		}
	}
	s.log.Infof("🚀 NSQ Server started, listening to %d channels", len(s.consumers))
	return nil
}

// Stop gracefully drains all channels and drops connections when Kratos shuts down
func (s *NSQServer) Stop(ctx context.Context) error {
	// 1. Tell all consumers to stop accepting new messages
	for _, consumer := range s.consumers {
		consumer.Stop()
	}

	// 2. Block until all currently processing messages are finished
	for _, consumer := range s.consumers {
		<-consumer.StopChan
	}
	s.log.Info("🛑 NSQ Server gracefully stopped")
	return nil
}
