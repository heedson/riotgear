package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
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

func (s *Server) GetPlayerID(ctx context.Context, pbReq *proto.PlayerIDReq) (*proto.PlayerID, error) {
	if pbReq.GetPlayerName() == "" {
		return nil, errors.New("player name cannot be empty")
	}

	rel := &url.URL{Path: fmt.Sprintf("/lol/summoner/v3/summoners/by-name/%s", pbReq.GetPlayerName())}

	u := s.baseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Riot-Token", s.riotAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	var data = make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.logger.Infoln(int(data["id"].(float64)))

	return &proto.PlayerID{
		PlayerId: int64(data["id"].(float64)),
	}, nil
}
