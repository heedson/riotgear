package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/heedson/riotgear/api"
	"github.com/heedson/riotgear/gear"
	"github.com/heedson/riotgear/proto"
)

var testPlayerNames = []string{
	"TestPlayer27",
	"ToxicTroll69",
}

var testPlayers = map[string]gear.PlayerData{
	strings.ToLower(testPlayerNames[0]): gear.PlayerData{
		ProfileIconID: 12345,
		Name:          testPlayerNames[0],
		SummonerLevel: 67.0,
		RevisionDate:  1521154178000.0,
		ID:            89012.0,
		AccountID:     34567.0,
	},
	strings.ToLower(testPlayerNames[1]): gear.PlayerData{
		ProfileIconID: 67890,
		Name:          testPlayerNames[1],
		SummonerLevel: 12.0,
		RevisionDate:  1521154178000.0,
		ID:            34567.0,
		AccountID:     89012.0,
	},
}

var testRankData = map[int][]gear.LeaguePositionData{
	int(testPlayers[strings.ToLower(testPlayerNames[0])].ID): []gear.LeaguePositionData{
		{
			Rank:             "II",
			QueueType:        "RANKED_SOLO_5x5",
			HotStreak:        false,
			MiniSeries:       gear.MiniSeriesData{},
			Wins:             84,
			Veteran:          true,
			Losses:           68,
			FreshBlood:       false,
			LeagueID:         "<UUID here>",
			PlayerOrTeamName: testPlayers["testplayer27"].Name,
			Inactive:         false,
			PlayerOrTeamID:   fmt.Sprintf("%d", int(testPlayers["testplayer27"].ID)),
			LeagueName:       "Singed's Riotshields",
			Tier:             "PLATINUM",
			LeaguePoints:     83,
		},
	},
	int(testPlayers[strings.ToLower(testPlayerNames[1])].ID): []gear.LeaguePositionData{
		{
			Rank:             "IV",
			QueueType:        "RANKED_SOLO_5x5",
			HotStreak:        false,
			MiniSeries:       gear.MiniSeriesData{},
			Wins:             28,
			Veteran:          false,
			Losses:           68,
			FreshBlood:       false,
			LeagueID:         "<UUID here>",
			PlayerOrTeamName: testPlayers["toxictroll69"].Name,
			Inactive:         false,
			PlayerOrTeamID:   fmt.Sprintf("%d", int(testPlayers["toxictroll69"].ID)),
			LeagueName:       "Singed's Riotshields",
			Tier:             "PLATINUM",
			LeaguePoints:     21,
		},
	},
}

type mockRiotHandler struct {
	apiKey string
	paths  map[string]string
}

func (r mockRiotHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	apiKeys, ok := req.Header["X-Riot-Token"]
	if !(ok && len(apiKeys) == 1 && apiKeys[0] == r.apiKey) {
		resp.WriteHeader(403)
		return
	}

	splitPath := strings.Split(req.URL.String(), "/")

	path := strings.Join(
		splitPath[:len(splitPath)-1],
		"/",
	)

	path += "/"

	var data []byte
	var err error

	switch path {
	case r.paths["playerdata"]:
		playerName := splitPath[len(splitPath)-1]
		player, ok := testPlayers[strings.ToLower(playerName)]
		if !ok {
			resp.WriteHeader(404)
			return
		}
		data, err = json.Marshal(player)
		if err != nil {
			panic(err)
		}
	case r.paths["rankdata"]:
		pathPlayerID := splitPath[len(splitPath)-1]
		playerID, err := strconv.Atoi(pathPlayerID)
		if err != nil {
			panic(err)
		}

		rankData, ok := testRankData[playerID]
		if !ok {
			resp.WriteHeader(404)
			return
		}

		data, err = json.Marshal(rankData)
		if err != nil {
			panic(err)
		}
	default:
		resp.WriteHeader(404)
		return
	}

	_, err = resp.Write([]byte(data))
	if err != nil {
		panic(err)
	}
}

func newMockRiotHandler(mockAPIKey string) http.Handler {
	return mockRiotHandler{
		apiKey: mockAPIKey,
		paths: map[string]string{
			"playerdata": "/lol/summoner/v3/summoners/by-name/",
			"rankdata":   "/lol/league/v3/positions/by-summoner/",
		},
	}
}

func TestGetPlayerRank(t *testing.T) {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: time.StampMilli,
		FullTimestamp:   true,
	}

	mockRiotServer := httptest.NewServer(newMockRiotHandler("myapikey"))
	defer mockRiotServer.Close()

	mockClient := mockRiotServer.Client()

	mockURL, err := url.Parse(mockRiotServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	mockRegionsToMockURLs := map[string]*url.URL{
		"test": mockURL,
	}

	s := api.NewServer(logger, mockClient, mockRegionsToMockURLs, "myapikey")

	for _, playerName := range testPlayerNames {
		t.Run(playerName, func(t *testing.T) {
			resp, err := s.GetPlayerRank(context.Background(), &proto.PlayerRankReq{
				RegionName: "test",
				PlayerName: testPlayers[strings.ToLower(playerName)].Name,
			})
			if err != nil {
				t.Fatal(err)
			}

			if resp.GetPlayerId() != int64(testPlayers[strings.ToLower(playerName)].ID) {
				t.Errorf("got %d; want %d", resp.GetPlayerId(), int64(testPlayers["testplayer27"].ID))
			}
		})
	}
}
