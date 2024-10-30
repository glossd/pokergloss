package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const MinSitAndGoBuyIn = 50
const DefaultLevelIncreaseTime = 5 * time.Minute

var ErrFullLobby = E("tournament lobby is full")
var ErrSitngoNotEnoughPlayers = E("not enough players")

type LobbySitAndGo struct {
	ID primitive.ObjectID `bson:"_id"`
	TournamentLobby

	Size int
	BettingLimit

	PlacesPaid        int
	BuyIn             int64
	LevelIncreaseTime time.Duration
	StartingChips     int64
	BigBlind          int64
	StartAt           int64

	Entries []*EntrySitAndGo

	// millis
	CreatedAt int64
	CreatedBy string

	TableID *primitive.ObjectID

	IsPrivate bool

	table *Table

	// Optimistic locking
	Version int

	// is set on db operation, you can't use it inside domain
	// created for fast filtering of empty and full tables
	PlayersCount int
}

type NewLobbySitAndGoParams struct {
	NewTableParams
	PlacesPaid        int
	BuyIn             int64
	StartAt           int64
	LevelIncreaseTime time.Duration
}

var MinLevelIncreaseMinDuration = time.Duration(1)

const MaxLevelIncreaseMinDuration = 60

func NewLobbySitAndGo(params NewLobbySitAndGoParams) (*LobbySitAndGo, error) {
	err := ValidateAndEnrichParams(&params.NewTableParams)
	if err != nil {
		return nil, err
	}

	if params.LevelIncreaseTime < MinLevelIncreaseMinDuration*time.Minute {
		return nil, E("level increase minutes minimum=%d", MinLevelIncreaseMinDuration)
	}
	if params.LevelIncreaseTime > MaxLevelIncreaseMinDuration*time.Minute {
		return nil, E("level increase minutes maximum=%d", MaxLevelIncreaseMinDuration)
	}
	if params.BettingLimit == "" {
		params.BettingLimit = NL
	}

	if params.StartAt != 0 {
		if params.StartAt < timeutil.Now() {
			return nil, E("start of sitngo must be in the future")
		}
		if params.StartAt > timeutil.NowAdd(goconf.Props.TableService.Cleaning.SitngoStartTimeout) {
			return nil, E("start of sitngo must be no further than %d minutes", int(goconf.Props.TableService.SitngoStartTimeout.Minutes()))
		}
		timeStartAt := timeutil.ToTime(params.StartAt)
		if timeStartAt.Second() != 0 || timeStartAt.Nanosecond() != 0 {
			return nil, E("start of sitngo must be at exact minute")
		}
	}

	if 5*params.NewTableParams.BigBlind > params.BuyIn {
		return nil, E("buy in must be at least 5xBB")
	}

	if params.BuyIn < MinSitAndGoBuyIn {
		return nil, E("minimum buy-in is %d", MinSitAndGoBuyIn)
	}

	maxPaidPlaces := maxPlacesPaid(params.Size)
	if params.PlacesPaid > maxPaidPlaces {
		return nil, E("maximum paid places for your table size is %d", params.Size/3)
	}

	if params.PlacesPaid <= 0 {
		return nil, E("minimum paid places is %d", 1)
	}

	tournament := TournamentLobby{
		Name:            params.Name,
		FeePercent:      GetTournamentFeePercent(),
		Status:          LobbyRegistering,
		DecisionTimeout: params.DecisionTimeout,
	}

	lobby := &LobbySitAndGo{
		ID:                primitive.NewObjectID(),
		TournamentLobby:   tournament,
		Entries:           []*EntrySitAndGo{},
		Size:              params.Size,
		BettingLimit:      params.BettingLimit,
		PlacesPaid:        params.PlacesPaid,
		BuyIn:             params.BuyIn,
		BigBlind:          params.BigBlind,
		LevelIncreaseTime: params.LevelIncreaseTime,
		StartAt:           params.StartAt,
		StartingChips:     params.BuyIn,
		IsPrivate:         params.IsPrivate,
		CreatedAt:         timeutil.Now(),
		CreatedBy:         params.Identity.UserId,
	}

	return lobby, nil
}

func (l *LobbySitAndGo) Register(iden authid.Identity, position int) error {
	if len(l.Entries) >= l.Size {
		return ErrFullLobby
	}

	for _, entry := range l.Entries {
		if entry.Position == position {
			return ErrSeatTaken
		}
		if entry.UserId == iden.UserId {
			return ErrAlreadySitting
		}
	}

	l.Entries = append(l.Entries, &EntrySitAndGo{Identity: iden, Position: position})
	l.Prizes = calcPrizes(l)

	if len(l.Entries) == l.Size {
		t, err := NewTableSitAndGo(l.tableParams(), l.ToSeats(), l.tournamentAttrs(), true)
		if err != nil {
			log.Errorf("Couldn't create table from sit&go tournament lobby: %s", err)
		}

		l.toStarted(t)

		return nil
	}

	return nil
}

func (l *LobbySitAndGo) toStarted(t *Table) {
	l.table = t
	l.TableID = &t.ID
	l.Status = LobbyStarted
}

func (l *LobbySitAndGo) StartAnyway() error {
	if len(l.Entries) < 2 {
		return ErrSitngoNotEnoughPlayers
	}
	if l.PlacesPaid > l.maxPlacesPaid() {
		l.PlacesPaid = l.maxPlacesPaid()
	}

	t, err := NewTableSitAndGo(l.tableParams(), l.ToSeats(), l.tournamentAttrs(), true)
	if err != nil {
		return err
	}
	l.toStarted(t)
	return nil
}

func (l *LobbySitAndGo) Unregister(iden authid.Identity, position int) error {
	if l.Status != LobbyRegistering {
		return E("tournament is started")
	}

	var foundUserIdx int
	var foundUserId string
	for i, entry := range l.Entries {
		if position == entry.Position {
			foundUserId = entry.UserId
			foundUserIdx = i
			break
		}
	}

	if len(foundUserId) == 0 {
		return E("seat is not taken")
	}

	if foundUserId != iden.UserId {
		return ErrPositionNotYours
	}

	l.Entries = append(l.Entries[:foundUserIdx], l.Entries[foundUserIdx+1:]...)
	l.Prizes = calcPrizes(l)

	return nil
}

func (l *LobbySitAndGo) maxPlacesPaid() int {
	return maxPlacesPaid(len(l.Entries))
}

func (l *LobbySitAndGo) tableParams() NewTableParams {
	return NewTableParams{
		Name:            l.Name,
		Size:            l.Size,
		BigBlind:        l.BigBlind,
		DecisionTimeout: l.DecisionTimeout,
		BettingLimit:    l.BettingLimit,
		Identity:        authid.Identity{UserId: l.CreatedBy},
	}
}

func (l *LobbySitAndGo) ToSeats() []*Seat {
	seats := make([]*Seat, 0, len(l.Entries))
	for _, e := range l.Entries {
		seats = append(seats, NewTournamentSeat(e.Position, e.Identity, l.StartingChips))
	}
	return seats
}

func (l *LobbySitAndGo) tournamentAttrs() TournamentAttributes {
	var attrs TournamentAttributes
	attrs.LobbyID = l.ID
	attrs.Name = l.Name
	attrs.StartAt = timeutil.Now()
	attrs.LevelIncreaseTime = l.LevelIncreaseTime
	attrs.LevelIncreaseAt = timeutil.NowAdd(l.LevelIncreaseTime)
	attrs.Prizes = l.Prizes
	attrs.NextSmallBlind = nextSmallBlind(l.BigBlind / 2)
	attrs.FeePercent = l.FeePercent
	attrs.BuyIn = l.BuyIn
	return attrs
}

func (l *LobbySitAndGo) GetTable() *Table {
	return l.table
}

func (l *LobbySitAndGo) PrizePool() int64 {
	poolFeeFree := l.BuyIn * int64(len(l.Entries))
	return poolFeeFree - computeTournamentFeeFrom(poolFeeFree, l.FeePercent)
}

func (l *LobbySitAndGo) ExpectedPrizePool() int64 {
	poolFeeFree := l.BuyIn * int64(l.Size)
	return poolFeeFree - computeTournamentFeeFrom(poolFeeFree, l.FeePercent)
}

func (l *LobbySitAndGo) GetPlacesPaid() int {
	return l.PlacesPaid
}
