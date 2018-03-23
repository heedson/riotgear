package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
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

// Server is the Riotgear server struct.
type Server struct {
	logger *logrus.Logger
	db     *sql.DB

	regionToServerURL map[string]*url.URL

	riotAPIKey string
	httpClient *http.Client

	nameRegexp *regexp.Regexp
}

// NewServer returns a new Riotgear server that utilises the given API key on
// the many region server URLs to satisfy queries.
func NewServer(logger *logrus.Logger, db *sql.DB, schemaSource io.Reader, client *http.Client, regionToServerURL map[string]*url.URL, riotAPIKey string) (*Server, error) {
	if db == nil {
		return nil, errors.New("DB cannot be nil")
	}

	s := &Server{
		logger:            logger,
		db:                db,
		regionToServerURL: regionToServerURL,
		riotAPIKey:        riotAPIKey,
		httpClient:        client,
		nameRegexp:        regexp.MustCompile(`^[0-9\p{L} _.]+$`),
	}

	if err := s.buildSchema(schemaSource); err != nil {
		return nil, errors.WithStack(err)
	}

	return s, nil
}

func (s *Server) buildSchema(schemaSource io.Reader) error {
	schema, err := ioutil.ReadAll(schemaSource)
	if err != nil {
		return errors.WithStack(err)
	}

	commands := strings.Split(string(schema), ";")

	tx, err := s.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}
	defer tx.Rollback()

	for _, cmd := range commands {
		_, err = tx.Exec(cmd)
		if err != nil {
			return err
		}
	}

	return errors.WithStack(tx.Commit())
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

// GetPlayerID returns the player ID when given the region and a player name.
// This is just an example use of Riot's League of Lengend's API.
func (s *Server) GetPlayerID(ctx context.Context, pbReq *proto.PlayerReq) (*proto.PlayerID, error) {
	serverURL, err := s.getServerURL(pbReq.GetRegionName())
	if err != nil {
		return nil, getGRPCError(err)
	}

	playerData, err := s.getPlayerData(serverURL, pbReq.GetPlayerName())
	if err != nil {
		return nil, getGRPCError(err)
	}

	s.logger.Infof("Request for %s on %s. Player ID %d.", playerData.Name, strings.ToUpper(pbReq.GetRegionName()), playerData.ID)

	return &proto.PlayerID{
		PlayerId: int64(playerData.ID),
	}, nil
}

// GetPlayerRank returns the rank stats for all queue types for a given player name
// on a given region.
func (s *Server) GetPlayerRank(ctx context.Context, pbReq *proto.PlayerReq) (*proto.PlayerRank, error) {
	serverURL, err := s.getServerURL(pbReq.GetRegionName())
	if err != nil {
		return nil, getGRPCError(err)
	}

	playerData, err := s.getPlayerData(serverURL, pbReq.GetPlayerName())
	if err != nil {
		return nil, getGRPCError(err)
	}

	s.logger.Infof("Request for %s on %s. Player ID %d.", playerData.Name, strings.ToUpper(pbReq.GetRegionName()), playerData.ID)

	rel := &url.URL{Path: fmt.Sprintf("/lol/league/v3/positions/by-summoner/%d", playerData.ID)}

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
		return nil, getGRPCError(shield.Errorf(shield.Internal, "%s. Player ID %d", resp.Status, playerData.ID))
	}

	var data gear.RankData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, getGRPCError(shield.Wrap(shield.Internal, err))
	}

	s.logger.Infof("%#v", data)

	pbResp := &proto.PlayerRank{
		LeaguePositions: make([]*proto.PlayerRank_LeaguePosition, 0, len(data)),
	}

	for _, lp := range data {
		leaguePosition := &proto.PlayerRank_LeaguePosition{
			Rank:      lp.Rank,
			QueueType: lp.QueueType,
			HotStreak: lp.HotStreak,
			MiniSeries: &proto.PlayerRank_LeaguePosition_MiniSeries{
				Wins:     int64(lp.MiniSeries.Wins),
				Losses:   int64(lp.MiniSeries.Losses),
				Target:   int64(lp.MiniSeries.Target),
				Progress: lp.MiniSeries.Progress,
			},
			Wins:             int64(lp.Wins),
			Veteran:          lp.Veteran,
			Losses:           int64(lp.Losses),
			FreshBlood:       lp.FreshBlood,
			LeagueId:         lp.LeagueID,
			PlayerOrTeamName: lp.PlayerOrTeamName,
			Inactive:         lp.Inactive,
			PlayerOrTeamId:   lp.PlayerOrTeamID,
			LeagueName:       lp.LeagueName,
			Tier:             lp.Tier,
			LeaguePoints:     int64(lp.LeaguePoints),
		}

		pbResp.LeaguePositions = append(pbResp.LeaguePositions, leaguePosition)
	}

	return pbResp, nil
}
