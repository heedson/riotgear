package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/heedson/riotgear/proto"
)

// Fix when dep isn't shit
//go:generate protoc -I../proto -I../thirdparty/googleapis/ --go_out=plugins=grpc:$GOPATH/src --grpc-gateway_out=logtostderr=true:$GOPATH/src ../proto/api.proto

type Server struct {
	logger *logrus.Logger

	baseURL    *url.URL
	riotAPIKey string

	httpClient *http.Client
}

func NewServer(logger *logrus.Logger, baseURL *url.URL, riotAPIKey string) *Server {
	return &Server{
		logger:     logger,
		baseURL:    baseURL,
		riotAPIKey: riotAPIKey,
		httpClient: http.DefaultClient,
	}
}

func (s *Server) Echo(ctx context.Context, req *proto.EchoMsg) (*proto.EchoMsg, error) {
	s.logger.Infoln(req.GetValue())

	http.Client

	return &proto.EchoMsg{
		Value: req.GetValue(),
	}, nil
}
