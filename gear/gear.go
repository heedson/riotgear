package gear

type RankData struct {
	Ranks []LeaguePositionData
}

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

type MiniSeriesData struct {
	Wins     int
	Losses   int
	Target   int
	Progress string
}
