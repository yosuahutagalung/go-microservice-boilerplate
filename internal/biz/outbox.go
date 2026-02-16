package biz

import "context"

type OutboxEvent struct {
	ID      string
	Topic   string
	Payload []byte
}

type OutboxRepo interface {
	GetPendingEvents(ctx context.Context, limit int32) ([]*OutboxEvent, error)
	MarkPublished(ctx context.Context, id string) error
}
