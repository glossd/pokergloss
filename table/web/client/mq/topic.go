package mq

const StartMultiTopicID = "pg.table.start-multi"

type StartTournamentEvent struct {
	StartAt int64 `json:"startAt"`
}

const StartSitngoTopicID = "pg.table.start-sitngo"

const TimeoutTopicID = "pg.table.timeout"

const MultiPlayersMovedTopicID = "pg.table.multi-players-moved"

type MultiPlayersMovedEvent struct {
	LobbyID     string `json:"lobbyId"`
	TableID     string `json:"tableId"`
	RebalanceAt int64  `json:"rebalanceAt"`
}
