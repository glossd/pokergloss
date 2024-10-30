package domain

import (
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"time"
)

var ErrNotAvailableInTournament = E("no available in tournament type game")

type TournamentAttributes struct {
	LobbyID           primitive.ObjectID
	Name              string
	StartAt           int64
	LevelIncreaseTime time.Duration
	Prizes            []TournamentPrize
	MarketPrize       *MarketPrize
	LevelIncreaseAt   int64
	NextSmallBlind    int64
	BuyIn             int64
	FeePercent        float64

	// for sitngo
	TournamentWinners []*TournamentWinner
}

func (t *Table) setPlayerTournamentInfo(p *Player, place int) (prize int64) {
	for _, tp := range t.Prizes {
		if tp.Place == place {
			prize = tp.Prize
		}
	}
	p.tournamentInfo = PlayerTournamentInfo{Place: place, Prize: prize}
	return prize
}

func maxPlacesPaid(users int) int {
	return int(math.Max(1, float64(users/3)))
}

func (t *Table) checkTimeAndIncreaseBlinds() {
	if timeutil.Now() >= t.LevelIncreaseAt {
		t.increaseBlinds()
		t.LevelIncreaseAt = timeutil.Add(t.LevelIncreaseAt, t.LevelIncreaseTime)
	}
}

func (t *Table) increaseBlinds() {
	nextSB := t.TournamentAttributes.NextSmallBlind
	t.BigBlind = nextSB * 2
	t.SmallBlind = nextSB
	t.TournamentAttributes.NextSmallBlind = nextSmallBlind(nextSB)
}

func (t *Table) GetNextSmallBlind() int64 {
	return t.NextSmallBlind
}

func (t *TournamentAttributes) Fee() int64 {
	return computeTournamentFeeFrom(t.BuyIn, t.FeePercent)
}

func nextSmallBlind(current int64) int64 {
	var i int64
	var tens int64
	for i = current; i >= 10; i = i / 10 {
		tens++
	}
	switch i {
	case 1:
		return i * powTen(tens) * 2
	case 2:
		return i * powTen(tens) / 2 * 3
	case 3:
		return i * powTen(tens) / 3 * 5
	case 4:
		return i * powTen(tens) / 4 * 5
	case 5:
		return i * powTen(tens) * 2
	default:
		log.Errorf("No nextSmallBlind is found for current=%d", current)
		return powTen(tens + 1)
	}
}

func powTen(tens int64) int64 {
	return int64(math.Pow10(int(tens)))
}
