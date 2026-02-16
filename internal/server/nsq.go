package server

import (
	"context"
	"fmt"

	"service_boilerplate/internal/conf"
	"service_boilerplate/internal/service"

	"github.com/nsqio/go-nsq"
)

// NSQServer manages all background message consumers
type NSQServer struct {
	consumers []*nsq.Consumer
	lookupd   string
}

func NewNSQServer(c *conf.Data, svc *service.GreeterService) (*NSQServer, error) {
	srv := &NSQServer{
		lookupd: c.Mq.LookupdAddr,
	}

	nsqConfig := nsq.NewConfig()

	// ==========================================
	// 1. Setup Greeting Events Consumer
	// ==========================================
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
		// Just a helpful warning in case someone forgot to add it to config.yaml
		fmt.Println("⚠️ WARNING: 'greeting_events' missing from config.yaml topics map")
	}

	// ==========================================
	// 2. Add Future Consumers Here (e.g., Loan Approval, Audit)
	// ==========================================
	/*
		if info, exists := c.Mq.Topics["loan_events"]; exists {
			// Initialize loan consumer, attach svc.HandleLoanEvent, append to srv.consumers
		}
	*/

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
	fmt.Printf("🚀 NSQ Server started, listening to %d channels\n", len(s.consumers))
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
	fmt.Println("🛑 NSQ Server gracefully stopped")
	return nil
}
