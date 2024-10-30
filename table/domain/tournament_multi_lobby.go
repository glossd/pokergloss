package domain

import (
	"fmt"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"strconv"
	"time"
)

const defaultTournamentBuyIn = 1000
const defaultTournamentBigBlind = 10
const defaultLevelIncreaseTime = 5 * time.Minute
const minBigBlind = 2

var ErrNotRegisteredForTournament = E("you are not registered for tournament")
var ErrTournamentStarted = E("tournament started")

type LobbyMulti struct {
	ID primitive.ObjectID `bson:"_id"`
	TournamentLobby

	StartAt  int64 // millis
	BuyIn    int64
	BigBlind int64

	TableSize         int
	LevelIncreaseTime time.Duration
	BettingLimit

	// If it's zero it means the tournament is a FreeRoll
	// BuyIn int64

	Users []authid.Identity

	TableIDs []primitive.ObjectID

	// supposed to be set manually in db
	RealPrizeDescription *string
	LastVideoID          *string
	MarketPrize          *MarketPrize

	CreatedAt int64

	Version int64

	// won't be saved to db
	tables []*Table
}

type NewLobbyMultiParams struct {
	StartAt   int64 // millis
	BuyIn     int64
	BigBlind  int64
	Name      string
	TableSize int
	BettingLimit
	LevelIncreaseTime time.Duration
	DecisionTimeout   time.Duration
}

type MarketPrize struct {
	ItemID       string
	NumberOfDays int
}

func NewLobbyMulti(params NewLobbyMultiParams) *LobbyMulti {
	tournament := TournamentLobby{
		Name:            params.Name,
		FeePercent:      GetTournamentFeePercent(),
		Status:          LobbyRegistering,
		DecisionTimeout: params.DecisionTimeout,
	}
	return &LobbyMulti{
		ID:                primitive.NewObjectID(),
		TournamentLobby:   tournament,
		CreatedAt:         timeutil.Now(),
		StartAt:           params.StartAt,
		BuyIn:             params.BuyIn,
		BigBlind:          params.BigBlind,
		TableSize:         params.TableSize,
		BettingLimit:      params.BettingLimit,
		LevelIncreaseTime: params.LevelIncreaseTime,
		Users:             nil,
	}
}

func NewMultiLobbyWithName(name string, bl BettingLimit, startAt time.Time) *LobbyMulti {
	params := defaultNewLobbyParams(name, startAt)
	params.BettingLimit = bl
	return NewLobbyMulti(params)
}

func NewMultiLobbyDaily(name string, bl BettingLimit, startAt time.Time, buyIn, bigBlind int64) *LobbyMulti {
	params := defaultNewLobbyParams(name, startAt)
	params.BettingLimit = bl
	params.BuyIn = buyIn
	params.BigBlind = bigBlind
	return NewLobbyMulti(params)
}

func NewFreerollWithNumber(startAt time.Time, number int) *LobbyMulti {
	name := fmt.Sprintf("Freeroll %d #%d", startAt.Day(), number)
	params := defaultNewLobbyParams(name, startAt)
	lobby := NewLobbyMulti(params)
	lobby.TournamentLobby.FeePercent = 0
	return lobby
}

func defaultNewLobbyParams(name string, startAt time.Time) NewLobbyMultiParams {
	return NewLobbyMultiParams{
		StartAt:           timeutil.Time(startAt),
		BuyIn:             defaultTournamentBuyIn,
		BigBlind:          defaultTournamentBigBlind,
		Name:              name,
		TableSize:         goconf.Props.TableService.Multi.TableSize,
		BettingLimit:      NL,
		LevelIncreaseTime: defaultLevelIncreaseTime,
		DecisionTimeout:   goconf.Props.TableService.Multi.DecisionTimeout,
	}
}

func (l *LobbyMulti) Register(iden authid.Identity) error {
	if l.StartAt < timeutil.Now() {
		return ErrTournamentStarted
	}

	for _, user := range l.Users {
		if user.UserId == iden.UserId {
			return ErrAlreadySitting
		}
	}

	l.Users = append(l.Users, iden)

	l.Prizes = calcPrizes(l)

	return nil
}

func (l *LobbyMulti) Unregister(iden authid.Identity) (int, error) {
	if l.StartAt < timeutil.Now() {
		return 0, ErrTournamentStarted
	}
	idx := -1
	for i, user := range l.Users {
		if user.UserId == iden.UserId {
			idx = i
		}
	}
	if idx < 0 {
		return idx, ErrNotRegisteredForTournament
	}

	// unordered remove
	l.Users[idx] = l.Users[len(l.Users)-1]
	l.Users = l.Users[:len(l.Users)-1]

	l.Prizes = calcPrizes(l)

	return idx, nil
}

// Users distribution:
// 7 - 3+1, 3
// 8 - 4, 4
// 13 - 4+1, 4, 4
// 31 - 5+1, 5, 5, 5, 5
func (l *LobbyMulti) Start() {
	if len(l.Users) < 2 {
		l.Status = LobbyFinished
		return
	}

	Algo.ShuffleUsers(l.Users)

	tablesNum := int(math.Ceil(float64(len(l.Users)) / float64(l.TableSize)))
	minUsersPerTable := int(math.Floor(float64(len(l.Users)) / float64(tablesNum)))
	remainder := len(l.Users) - minUsersPerTable*tablesNum
	for i := 0; i < tablesNum; i++ {
		from := i * minUsersPerTable
		to := (i + 1) * minUsersPerTable
		seats := l.IdensToSeats(l.Users[from:to])
		if remainder > 0 {
			seats = append(seats, l.IdenToSeat(minUsersPerTable, l.Users[len(l.Users)-remainder]))
			remainder--
		}
		table, err := NewTableMulti(l, l.tableParams(i+1), seats)
		if err != nil {
			log.Errorf("Failed to start multi lobby: %s", err)
			continue
		}
		if tablesNum == 1 {
			table.IsLast = true
		}
		l.TableIDs = append(l.TableIDs, table.ID)
		l.tables = append(l.tables, table)
	}
	l.Status = LobbyStarted
}

func (l *LobbyMulti) IdensToSeats(idens []authid.Identity) []*Seat {
	var result []*Seat
	for i, iden := range idens {
		result = append(result, l.IdenToSeat(i, iden))
	}
	return result
}

func (l *LobbyMulti) IdenToSeat(position int, iden authid.Identity) *Seat {
	return NewTournamentSeat(position, iden, l.BuyIn)
}

func (l *LobbyMulti) GetTables() []*Table {
	return l.tables
}

func (l *LobbyMulti) tableParams(tableNum int) NewTableParams {
	return NewTableParams{
		Name:            strconv.Itoa(tableNum),
		Size:            l.TableSize,
		BigBlind:        l.BigBlind,
		BettingLimit:    l.BettingLimit,
		DecisionTimeout: l.DecisionTimeout,
		Identity:        authid.Identity{UserId: "Automatic"},
	}
}

func (l *LobbyMulti) tournamentAttrs() TournamentAttributes {
	var attrs TournamentAttributes
	attrs.LobbyID = l.ID
	attrs.Name = l.Name
	attrs.StartAt = l.StartAt
	attrs.LevelIncreaseTime = l.LevelIncreaseTime
	attrs.LevelIncreaseAt = timeutil.NowAdd(l.LevelIncreaseTime)
	attrs.Prizes = l.Prizes
	attrs.MarketPrize = l.MarketPrize
	attrs.NextSmallBlind = nextSmallBlind(minBigBlind / 2)
	attrs.FeePercent = l.FeePercent
	attrs.BuyIn = l.BuyIn
	return attrs
}

func (l *LobbyMulti) GetPlacesPaid() int {
	return maxPlacesPaid(len(l.Users))
}

func (l *LobbyMulti) PrizePool() int64 {
	var poolFeeFree = l.BuyIn * int64(len(l.Users))
	return poolFeeFree - computeTournamentFeeFrom(poolFeeFree, l.FeePercent)
}

func (l *LobbyMulti) GetTableIDsAsStr() (a []string) {
	for _, id := range l.TableIDs {
		a = append(a, id.Hex())
	}
	return
}
