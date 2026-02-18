package data

import (
	"context"
	"time"

	"service_boilerplate/internal/biz"

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
	res := biz.Greeter{ID: g.ID, Hello: g.Hello, CreatedAt: time.Now()}
	return &res, nil
}
