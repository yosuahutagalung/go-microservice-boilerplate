package biz

import (
	"context"
	"time"

	v1 "github.com/yosuahutagalung/go-microservice-boilerplate-schema-registry/proto/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type GreeterEventPublisher interface {
	PublishGreetingSaid(ctx context.Context, g *Greeter) error
}

type Greeter struct {
	ID        string
	Hello     string
	CreatedAt time.Time
}

type GreeterRepo interface {
	Create(context.Context, *Greeter) (*Greeter, error)
}

type GreeterUsecase struct {
	repo GreeterRepo
	pub  GreeterEventPublisher
	log  *log.Helper
}

func NewGreeterUsecase(repo GreeterRepo, pub GreeterEventPublisher, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{
		repo: repo,
		pub:  pub,
		log:  log.NewHelper(logger),
	}
}

func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.ID)

	saved, err := uc.repo.Create(ctx, g)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (uc *GreeterUsecase) CreateGreeterAsync(ctx context.Context, g *Greeter) error {
	uc.log.WithContext(ctx).Infof("CreateGreeterEvent: %v", g.ID)
	err := uc.pub.PublishGreetingSaid(ctx, g)
	return err
}
