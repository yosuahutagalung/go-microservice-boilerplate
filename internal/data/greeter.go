package data

import (
	"context"

	"service_boilerplate/internal/biz"
	"service_boilerplate/internal/data/db"

	"github.com/go-kratos/kratos/v2/log"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) Create(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	i, err := r.data.query.CreateGreeter(ctx, db.CreateGreeterParams{
		ID:    g.ID,
		Hello: g.Hello,
	})
	if err != nil {
		return nil, err
	}

	res := biz.Greeter{ID: i.ID, Hello: i.Hello, CreatedAt: i.CreatedAt.Time}
	return &res, nil
}
