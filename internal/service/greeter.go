package service

import (
	"context"
	"encoding/json"

	v1 "service_boilerplate/api/helloworld/v1"
	"service_boilerplate/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nsqio/go-nsq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc  *biz.GreeterUsecase
	log *log.Helper
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{ID: in.Id, Hello: in.Hello})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Id: g.ID, Hello: g.Hello, CreatedAt: timestamppb.New(g.CreatedAt)}, nil
}

// HandleGreetingEvent processes incoming NSQ messages
func (s *GreeterService) HandleGreetingEvent(m *nsq.Message) error {
	var payload biz.Greeter
	if err := json.Unmarshal(m.Body, &payload); err != nil {
		return err // Returning error tells NSQ to requeue
	}
	_, err := s.uc.CreateGreeter(context.Background(), &payload)

	return err // Returning nil acknowledges the message is done
}

func (s *GreeterService) SayHelloAsync(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReplyAsync, error) {
	err := s.uc.CreateGreeterAsync(ctx, &biz.Greeter{ID: in.Id, Hello: in.Hello})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReplyAsync{Message: "Say hello event!"}, nil
}
