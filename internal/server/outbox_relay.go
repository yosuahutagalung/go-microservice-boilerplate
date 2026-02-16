package server

import (
	"context"
	"fmt"
	"time"

	"service_boilerplate/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nsqio/go-nsq"
)

type OutboxRelayServer struct {
	repo     biz.OutboxRepo
	producer *nsq.Producer // We inject the NSQ Producer directly here
	log      *log.Helper

	// Context controls the background loop lifecycle
	ctx    context.Context
	cancel context.CancelFunc
}

func NewOutboxRelayServer(repo biz.OutboxRepo, producer *nsq.Producer, logger log.Logger) *OutboxRelayServer {
	if producer == nil {
		log.NewHelper(logger).Errorf("❌ CRITICAL: NewOutboxRelayServer received NIL producer!")
	} else {
		log.NewHelper(logger).Infof("✅ NewOutboxRelayServer initialized with producer: %v", producer)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &OutboxRelayServer{
		repo:     repo,
		producer: producer,
		log:      log.NewHelper(logger),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start is called by Kratos automatically
func (s *OutboxRelayServer) Start(ctx context.Context) error {
	s.log.Infof("🚀 Starting Outbox Relay Worker... Producer: %v", s.producer)

	// Launch the infinite loop in a background goroutine
	go s.runLoop()
	return nil
}

// Stop is called by Kratos during graceful shutdown (e.g., pressing Ctrl+C)
func (s *OutboxRelayServer) Stop(ctx context.Context) error {
	s.log.Info("🛑 Stopping Outbox Relay Worker...")
	s.cancel() // This tells the loop to exit
	return nil
}

// The actual infinite loop
func (s *OutboxRelayServer) runLoop() {
	// Poll the database every 1 second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			// The server is shutting down, exit the loop cleanly
			return

		case <-ticker.C:
			// Time to check the database!
			s.processPendingEvents()
		}
	}
}

func (s *OutboxRelayServer) processPendingEvents() {
	// 1. Fetch up to 100 pending events (Using SKIP LOCKED so multiple workers don't collide)
	events, err := s.repo.GetPendingEvents(s.ctx, 100)
	if err != nil {
		s.log.Errorf("failed to fetch outbox events: %v", err)
		return
	}

	if len(events) == 0 {
		return // Nothing to do
	}

	// 2. Loop through and publish each one
	for _, event := range events {
		if s.producer == nil {
			s.log.Errorf("❌ ERROR: Producer is nil when trying to publish event %s", event.ID)
			return
		}
		// Publish the raw bytes directly to the exact topic requested
		err := s.producer.Publish(event.Topic, event.Payload)
		if err != nil {
			s.log.Errorf("failed to publish event %s: %v", event.ID, err)
			continue // If NSQ is down, skip this event. We will try again in 1 second.
		}

		// 3. Mark as published in the database
		err = s.repo.MarkPublished(s.ctx, event.ID)
		if err != nil {
			// Edge case: Message sent to NSQ, but DB update failed.
			// It will be re-sent next loop. (Downstream consumers must be idempotent!)
			s.log.Errorf("failed to mark event %s as published: %v", event.ID, err)
		} else {
			fmt.Printf("✅ Relayed outbox event: %s\n", event.ID)
		}
	}
}
