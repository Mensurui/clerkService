package service

import (
	"context"

	protos "github.com/Mensurui/clerkService/proto/golang/service"
	"github.com/hashicorp/go-hclog"
)

type Service struct {
	protos.UnimplementedServiceServer
	l hclog.Logger
}

func NewService(l hclog.Logger) *Service {
	return &Service{l: l}
}

func (s *Service) GetList(ctx context.Context, req *protos.GetListRequest) (*protos.GetListResponse, error) {
	return &protos.GetListResponse{
		List: []string{"one", "two", "three"},
	}, nil
}
