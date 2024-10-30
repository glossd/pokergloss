package domain

import (
	"fmt"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"time"
)

type ThemeID string

const (
	Paradise ThemeID = "paradise"
	Earth    ThemeID = "spirit-world"
	Hell     ThemeID = "hell"
)

func (t ThemeID) Name() string {
	switch t {
	case Paradise:
		return "Paradise"
	case Earth:
		return "Earth"
	case Hell:
		return "Hell"
	default:
		return "Table"
	}
}

type NewSurvivalTableParams struct {
	User              User
	Name              string
	DecisionTimeSec   int64
	BigBlind          int64
	Bots              []Bot
	ThemeID           ThemeID
	LevelIncreaseTime time.Duration
	SurvivalLevel     int64
}

type User struct {
	Iden  authid.Identity
	Stack int64
}

type Bot struct {
	Name    string
	Picture string
	Stack   int64
}

func NewTableSurvival(params NewSurvivalTableParams) (*Table, error) {
	tableParams := NewTableParams{
		Name:            "tmp",
		Size:            len(params.Bots) + 1,
		BigBlind:        params.BigBlind,
		DecisionTimeout: time.Duration(params.DecisionTimeSec) * time.Second,
		BettingLimit:    NL,
		IsPrivate:       true,
		Identity:        params.User.Iden,
	}

	var seats []*Seat
	seats = append(seats, NewTournamentSeat(0, params.User.Iden, params.User.Stack))
	for i, bot := range params.Bots {
		position := i + 1
		iden := authid.Identity{
			UserId:   fmt.Sprintf("%s_%d", bot.Name, position),
			Username: bot.Name,
			Picture:  bot.Picture,
		}
		seats = append(seats, NewTournamentSeat(position, iden, bot.Stack))
	}

	attrs := TournamentAttributes{
		StartAt:           timeutil.Now(),
		LevelIncreaseTime: params.LevelIncreaseTime,
		Prizes:            nil,
		MarketPrize:       nil,
		LevelIncreaseAt:   timeutil.NowAdd(params.LevelIncreaseTime),
		NextSmallBlind:    nextSmallBlind(params.BigBlind / 2),
		TournamentWinners: nil,
	}

	t, err := NewTableSitAndGo(tableParams, seats, attrs, false)
	if err != nil {
		return nil, err
	}

	t.Name = params.Name
	t.IsSurvival = true
	t.ThemeID = params.ThemeID
	t.SurvivalLevel = params.SurvivalLevel

	return t, nil
}

func (t *Table) IsSurvivalUserLeft() bool {
	return t.IsSurvival && t.IsSeatFree(0)
}

func (t *Table) MakeBotAction(position int, a Action) error {
	player, err := t.GetPlayer(position)
	if err != nil {
		return err
	}
	return t.makeActionWithValidation(player, a)
}

func (t *Table) SitBotBack(position int) error {
	p, err := t.GetPlayer(position)
	if err != nil {
		return err
	}

	if p.Status != PlayerSittingOut {
		return ErrNoSitBackInSittingOut
	}

	if p.Stack == 0 {
		return E("you can't sit back, you don't have chips on table")
	}

	p.Status = PlayerReady
	return nil
}
