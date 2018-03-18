package gear

import "encoding/json"

// PlayerData is the data object that contains data tied to the player.
type PlayerData struct {
	ProfileIconID int
	Name          string
	SummonerLevel int
	RevisionDate  int
	ID            int
	AccountID     int
}

func (p *PlayerData) UnmarshalJSON(data []byte) error {
	type aux struct {
		ProfileIconID int
		Name          string
		SummonerLevel float64
		RevisionDate  float64
		ID            float64
		AccountID     float64
	}

	var a aux

	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	p.ProfileIconID = a.ProfileIconID
	p.Name = a.Name
	p.SummonerLevel = int(a.SummonerLevel)
	p.RevisionDate = int(a.RevisionDate)
	p.ID = int(a.ID)
	p.AccountID = int(a.AccountID)

	return nil
}

// RankData is a slice of LeaguePositionData to represent a single player's whole
// list of ranks.
type RankData []LeaguePositionData

// LeaguePositionData is the complete data object for an individual ranked queue for
// a player or team.
type LeaguePositionData struct {
	Rank             string
	QueueType        string
	HotStreak        bool
	MiniSeries       MiniSeriesData
	Wins             int
	Veteran          bool
	Losses           int
	FreshBlood       bool
	LeagueID         string
	PlayerOrTeamName string
	Inactive         bool
	PlayerOrTeamID   string
	LeagueName       string
	Tier             string
	LeaguePoints     int
}

// MiniSeriesData is a data object for a mini series that the player or team may
// be in in a ranked queue.
type MiniSeriesData struct {
	Wins     int
	Losses   int
	Target   int
	Progress string
}
