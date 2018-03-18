package api_test

import (
	"context"
	"encoding/json"
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

type gearData struct {
	playerData gear.PlayerData
	rankData   gear.RankData
}

var testPlayers = map[string]gearData{
	"testplayer27": gearData{
		playerData: gear.PlayerData{
			ProfileIconID: 12345,
			Name:          "TestPlayer27",
			SummonerLevel: 67,
			RevisionDate:  1521154178000,
			ID:            89012,
			AccountID:     34567,
		},
		rankData: gear.RankData{
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
				PlayerOrTeamName: "TestPlayer27",
				Inactive:         false,
				PlayerOrTeamID:   strconv.Itoa(89012),
				LeagueName:       "Singed's Riotshields",
				Tier:             "PLATINUM",
				LeaguePoints:     83,
			},
		},
	},
	"toxictroll69": gearData{
		playerData: gear.PlayerData{
			ProfileIconID: 67890,
			Name:          "ToxicTroll69",
			SummonerLevel: 12,
			RevisionDate:  1521154178000,
			ID:            34567,
			AccountID:     89012,
		},
		rankData: gear.RankData{
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
				PlayerOrTeamName: "ToxicTroll69",
				Inactive:         false,
				PlayerOrTeamID:   strconv.Itoa(34567),
				LeagueName:       "Singed's Riotshields",
				Tier:             "PLATINUM",
				LeaguePoints:     21,
			},
		},
	},
}

func getTestPlayerByID(id int) (gearData, bool) {
	for _, player := range testPlayers {
		if id == int(player.playerData.ID) {
			return player, true
		}
	}

	return gearData{}, false
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

		data, err = json.Marshal(player.playerData)
		if err != nil {
			panic(err)
		}
	case r.paths["rankdata"]:
		pathPlayerID := splitPath[len(splitPath)-1]
		playerID, err := strconv.Atoi(pathPlayerID)
		if err != nil {
			panic(err)
		}

		player, ok := getTestPlayerByID(playerID)
		if !ok {
			resp.WriteHeader(404)
			return
		}

		data, err = json.Marshal(player.rankData)
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

	for name, player := range testPlayers {
		t.Run(name, func(t *testing.T) {
			resp, err := s.GetPlayerRank(context.Background(), &proto.PlayerReq{
				RegionName: "test",
				PlayerName: player.playerData.Name,
			})
			if err != nil {
				t.Fatal(err)
			}

			if len(resp.GetLeaguePositions()) != len(player.rankData) {
				t.Errorf("got %d league positions; want %d league positions", len(resp.GetLeaguePositions()), len(player.rankData))
			}
		})
	}
}
