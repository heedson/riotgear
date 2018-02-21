package api

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/heedson/riotgear/proto"
)

// Fix when dep isn't shit
//go:generate protoc -I../proto -I../thirdparty/googleapis/ --go_out=plugins=grpc:$GOPATH/src --grpc-gateway_out=logtostderr=true:$GOPATH/src ../proto/api.proto

type Server struct {
	logger *logrus.Logger
}

func NewServer(logger *logrus.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Echo(ctx context.Context, req *proto.EchoMsg) (*proto.EchoMsg, error) {
	s.logger.Infoln(req.GetValue())

	return &proto.EchoMsg{
		Value: req.GetValue(),
	}, nil
}
