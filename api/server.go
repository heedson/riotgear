package api

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
)

//go:generate protoc -I../proto -I../vendor/github.com/googleapis/googleapis/ --go_out=plugins=grpc:$GOPATH/src --grpc-gateway_out=logtostderr=true:../proto/ --swagger_out=logtostderr=true:../proto/ ../proto/api.proto

type Server struct {
	logger *logrus.Logger
}

func NewServer(logger *logrus.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Echo(ctx context.Context, req *proto.EchoRequest) (*empty.Empty, error) {
	fmt.Println(req.GetValue())
	return nil, nil
}
