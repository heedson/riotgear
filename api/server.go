package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

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
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *Server) Echo(ctx context.Context, req *proto.EchoMsg) (*proto.EchoMsg, error) {
	rel := &url.URL{Path: "/lol/league/v3/positions/by-summoner/heedson"}

	u := s.baseURL.ResolveReference(rel)
	re, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	re.Header.Set("Accept", "application/json")
	re.Header.Set("X-Riot-Token", s.riotAPIKey)

	resp, err := s.httpClient.Do(re)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data = make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	s.logger.Infof("%#v\n", data)
	s.logger.Infoln(req.GetValue())

	return &proto.EchoMsg{
		Value: req.GetValue(),
	}, nil
}
