package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/heedson/riotgear/gear"
	"github.com/heedson/riotgear/proto"
)

func getGRPCError(err error) error {
	switch err.(type) {
	case gear.ErrInternal:
		return status.Error(codes.Internal, err.Error())
	case gear.ErrInvalidArgument:
		return status.Error(codes.InvalidArgument, err.Error())
	case nil:
		return nil
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

type Server struct {
	logger *logrus.Logger

	regionToServerURL map[string]*url.URL

	riotAPIKey string
	httpClient *http.Client

	nameRegexp *regexp.Regexp
}

func NewServer(logger *logrus.Logger, regionToServerURL map[string]*url.URL, riotAPIKey string) *Server {
	return &Server{
		logger:            logger,
		regionToServerURL: regionToServerURL,
		riotAPIKey:        riotAPIKey,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		nameRegexp: regexp.MustCompile(`^[0-9\p{L} _.]+$`),
	}
}

func (s *Server) getPlayerID(ctx context.Context, regionName, playerName string) (int, error) {
	serverURL, ok := s.regionToServerURL[strings.ToLower(regionName)]
	if !ok {
		return 0, gear.Errorf(gear.InvalidArgument, "%q is not a valid region name", regionName)
	}

	if ok := s.nameRegexp.Match([]byte(playerName)); !ok {
		return 0, gear.Errorf(gear.InvalidArgument, "%q is not a valid player name", playerName)
	}

	rel := &url.URL{Path: "/lol/summoner/v3/summoners/by-name/" + playerName}

	u := serverURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, gear.Wrap(gear.Internal, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Riot-Token", s.riotAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, gear.Wrap(gear.Internal, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, gear.Error(gear.Internal, resp.Status)
	}

	var data = make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, gear.Wrap(gear.Internal, err)
	}

	playerID := int(data["id"].(float64))

	s.logger.Infof("%s's player ID: %d", data["name"], playerID)

	return playerID, nil
}

func (s *Server) GetPlayerID(ctx context.Context, pbReq *proto.PlayerIDReq) (*proto.PlayerID, error) {
	playerID, err := s.getPlayerID(ctx, pbReq.GetRegionName(), pbReq.GetPlayerName())
	if err != nil {
		return nil, getGRPCError(err)
	}

	return &proto.PlayerID{
		PlayerId: int64(playerID),
	}, nil
}
