package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/heedson/riotgear/gear"
	"github.com/heedson/riotgear/proto"
	"github.com/heedson/riotgear/shield"
)

func getGRPCError(err error) error {
	switch err.(type) {
	case shield.ErrInternal:
		return status.Error(codes.Internal, err.Error())
	case shield.ErrInvalidArgument:
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

func NewServer(logger *logrus.Logger, client *http.Client, regionToServerURL map[string]*url.URL, riotAPIKey string) *Server {
	return &Server{
		logger:            logger,
		regionToServerURL: regionToServerURL,
		riotAPIKey:        riotAPIKey,
		httpClient:        client,
		nameRegexp:        regexp.MustCompile(`^[0-9\p{L} _.]+$`),
	}
}

func (s *Server) getServerURL(regionName string) (*url.URL, error) {
	serverURL, ok := s.regionToServerURL[strings.ToLower(regionName)]
	if !ok {
		return nil, shield.Errorf(shield.InvalidArgument, "%q is not a valid region name", regionName)
	}

	return serverURL, nil
}

func (s *Server) getPlayerData(serverURL *url.URL, playerName string) (*gear.PlayerData, error) {
	if ok := s.nameRegexp.Match([]byte(playerName)); !ok {
		return nil, shield.Errorf(shield.InvalidArgument, "%q is not a valid player name", playerName)
	}

	rel := &url.URL{Path: "/lol/summoner/v3/summoners/by-name/" + playerName}

	u := serverURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, shield.Wrap(shield.Internal, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Riot-Token", s.riotAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, shield.Wrap(shield.Internal, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, shield.Errorf(shield.Internal, "%s. Player name %q", resp.Status, playerName)
	}

	var data gear.PlayerData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, shield.Wrap(shield.Internal, err)
	}

	return &data, nil
}

func (s *Server) GetPlayerID(ctx context.Context, pbReq *proto.PlayerIDReq) (*proto.PlayerID, error) {
	serverURL, err := s.getServerURL(pbReq.GetRegionName())
	if err != nil {
		return nil, getGRPCError(err)
	}

	playerData, err := s.getPlayerData(serverURL, pbReq.GetPlayerName())
	if err != nil {
		return nil, getGRPCError(err)
	}

	s.logger.Infof("Request for %s on %s. Player ID %d.", playerData.Name, strings.ToUpper(pbReq.GetRegionName()), int(playerData.ID))

	return &proto.PlayerID{
		PlayerId: int64(playerData.ID),
	}, nil
}

func (s *Server) GetPlayerRank(ctx context.Context, pbReq *proto.PlayerRankReq) (*proto.PlayerID, error) {
	serverURL, err := s.getServerURL(pbReq.GetRegionName())
	if err != nil {
		return nil, getGRPCError(err)
	}

	playerData, err := s.getPlayerData(serverURL, pbReq.GetPlayerName())
	if err != nil {
		return nil, getGRPCError(err)
	}

	s.logger.Infof("Request for %s on %s. Player ID %d.", playerData.Name, strings.ToUpper(pbReq.GetRegionName()), int(playerData.ID))

	rel := &url.URL{Path: fmt.Sprintf("/lol/league/v3/positions/by-summoner/%d", int(playerData.ID))}

	u := serverURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, getGRPCError(shield.Wrap(shield.Internal, err))
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Riot-Token", s.riotAPIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, getGRPCError(shield.Wrap(shield.Internal, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, getGRPCError(shield.Errorf(shield.Internal, "%s. Player ID %d", resp.Status, int(playerData.ID)))
	}

	//var rankData gear.RankData
	var data []gear.LeaguePositionData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, getGRPCError(shield.Wrap(shield.Internal, err))
	}

	s.logger.Infof("%#v", data)

	return &proto.PlayerID{
		PlayerId: int64(playerData.ID),
	}, nil
}
