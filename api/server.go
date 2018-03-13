package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/heedson/riotgear/proto"
)

var httpStatusToGRPCStatus = map[int]codes.Code{
	// 400 Bad Request - INTERNAL.
	http.StatusBadRequest: codes.Internal,
	// 401 Unauthorized  - UNAUTHENTICATED.
	http.StatusUnauthorized: codes.Unauthenticated,
	// 403 Forbidden - PERMISSION_DENIED.
	http.StatusForbidden: codes.PermissionDenied,
	// 404 Not Found - UNIMPLEMENTED.
	http.StatusNotFound: codes.Unimplemented,
	// 429 Too Many Requests - UNAVAILABLE.
	http.StatusTooManyRequests: codes.Unavailable,
	// 502 Bad Gateway - UNAVAILABLE.
	http.StatusBadGateway: codes.Unavailable,
	// 503 Service Unavailable - UNAVAILABLE.
	http.StatusServiceUnavailable: codes.Unavailable,
	// 504 Gateway timeout - UNAVAILABLE.
	http.StatusGatewayTimeout: codes.Unavailable,
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

func (s *Server) GetPlayerID(ctx context.Context, pbReq *proto.PlayerIDReq) (*proto.PlayerID, error) {
	if ok := s.nameRegexp.Match([]byte(pbReq.GetPlayerName())); !ok {
		return nil, status.Errorf(codes.InvalidArgument, "%q is not a valid player name", pbReq.GetPlayerName())
	}

	rel := &url.URL{Path: fmt.Sprintf("/lol/summoner/v3/summoners/by-name/%s", pbReq.GetPlayerName())}

	serverURL, ok := s.regionToServerURL[strings.ToLower(pbReq.GetRegion())]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "%q is not a valid server region", pbReq.GetRegion())
	}

	u := serverURL.ResolveReference(rel)

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

	if resp.StatusCode != http.StatusOK {
		grpcCode, ok := httpStatusToGRPCStatus[resp.StatusCode]
		if !ok {
			grpcCode = codes.Unknown
		}

		return nil, status.Error(grpcCode, resp.Status)
	}

	var data = make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	playerID := int(data["id"].(float64))

	s.logger.Infof("%s's player ID: %d\n", data["name"], playerID)

	return &proto.PlayerID{
		PlayerId: int64(playerID),
	}, nil
}
