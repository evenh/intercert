package server

import (
	"context"
	"github.com/evenh/intercert/api"
	"google.golang.org/grpc"
)

type Actuator struct{}

func (Actuator) Ping(context.Context, *api.PingMessage) (*api.PingMessage, error) {
	panic("implement me")
}

type Bar struct{}

func (Bar) Ping(ctx context.Context, in *api.PingMessage, opts ...grpc.CallOption) (*api.PingMessage, error) {
	panic("implement me")
}
