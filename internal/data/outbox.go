package data

import (
	"context"
	"service_boilerplate/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type outboxRepo struct {
	data *Data
	log  *log.Helper
}

func NewOutboxRepo(data *Data, logger log.Logger) biz.OutboxRepo {
	return &outboxRepo{data: data, log: log.NewHelper(logger)}
}

func (r *outboxRepo) GetPendingEvents(ctx context.Context, limit int32) ([]*biz.OutboxEvent, error) {
	rows, err := r.data.query.GetPendingOutboxEvents(ctx, limit)
	if err != nil {
		return nil, err
	}

	events := make([]*biz.OutboxEvent, 0, len(rows))

	for _, row := range rows {
		events = append(events, &biz.OutboxEvent{
			ID:      row.ID.String(),
			Topic:   row.Topic,
			Payload: row.Payload,
		})
	}

	return events, nil
}

func (r *outboxRepo) MarkPublished(ctx context.Context, id string) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = r.data.query.MarkOutboxEventPublished(ctx, uuidId)
	return err
}
