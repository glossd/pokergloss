package domain

import (
	"math"
	"time"
)

type LobbyStatus string
const (
	LobbyRegistering LobbyStatus = "registering"
	LobbyStarted     LobbyStatus = "started"
	LobbyFinished    LobbyStatus = "finished"
)

type TournamentPrize struct {
	Place int `json:"place"`
	Prize int64 `json:"prize"`
}

type TournamentLobby struct {
	Name string
	// entry fee is being taken from each player
	// or you could say it's taken from prize pool, same thing.
	// 1 - takes 100%, leaving players with nothing.
	// 0 - means no fee, basically making it a freeroll.
	FeePercent float64
	Status LobbyStatus
	DecisionTimeout time.Duration
	Prizes []TournamentPrize
}

type PrizeComputer interface {
	GetPlacesPaid() int
	PrizePool() int64
}

func calcPrizes(l PrizeComputer) []TournamentPrize {
	var prizes []TournamentPrize
	paid := l.GetPlacesPaid()
	for i := 0; i < paid; i++ {
		place := i + 1
		prize := ComputeTournamentPrize(l, place)
		prizes = append(prizes, TournamentPrize{Place: place, Prize: prize})
	}
	return prizes
}

// 1) 100% / 1 # 100%
// 2) 100% / 2, diff = 50% / 2 # 50% + diff, 50% - diff
// 3) diff = 33% / 2 # 33% + diff, 33%, 33% - diff
// 4) diff = 25% / 3 # 25% + 2*diff, 25% + diff, 25%- diff, 25% - 2*diff
// 5) diff = 20% / 4 # 20% + 2*diff, 20% + diff, 20%, 20% - diff, 20% - 3*diff
func ComputeTournamentPrize(l PrizeComputer, place int) int64 {
	placesPaid := l.GetPlacesPaid()
	if place > placesPaid {
		return 0
	}
	prizePool := l.PrizePool()
	if placesPaid == 1 {
		return prizePool
	}

	share := prizePool / int64(placesPaid)
	diff := share / int64(math.Max(2, float64(placesPaid) - 1))
	var factor = int64(placesPaid/2 + 1 - place)

	if placesPaid% 2 == 0 {
		if place > placesPaid/ 2 {
			factor--
		}
	}

	return share + factor * diff
}
