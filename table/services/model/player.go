package model

import (
	"github.com/glossd/pokergloss/table/domain"
)

var False = false
var True = true
var EmptyBlind = domain.Blind("")
var EmptyAction = domain.ActionType("")
var Zero int64 = 0
var NegativeOne = -1

type PlayerMapper func(*domain.Player) *Player

type PlayerAfterMapper func(*Player)

type Player struct {
	UserId   *string `json:"userId,omitempty"`
	Username *string `json:"username,omitempty"`
	Picture  *string `json:"picture,omitempty"`
	Position *int    `json:"position,omitempty"`

	Stack  *int64               `json:"stack,omitempty"`
	Status *domain.PlayerStatus `json:"status,omitempty" enums:"sitting,ready,playing,sittingOut"`
	Cards  *[]Card              `json:"cards,omitempty"`
	Blind  *domain.Blind        `json:"blind,omitempty" enums:"bigBlind,smallBlind,dealer,dealerSmallBlind"`

	TotalRoundBet *int64 `json:"totalRoundBet,omitempty"`

	LastGameBet    *int64             `json:"lastGameBet,omitempty"`
	LastGameAction *domain.ActionType `json:"lastGameAction,omitempty" enums:"check,bet,fold,call,raise,allIn"`

	ShowDownAction *domain.ShowDownActionType `json:"showDownAction,omitempty" enums:"show,muck,showLeft,showRight"`

	TimeoutAt *int64 `json:"timeoutAt,omitempty"`

	Intent *Intent `json:"intent"`

	Level           *int64  `json:"level,omitempty"`
	BankBalance     *int64  `json:"bankBalance,omitempty"`
	BankRank        *int64  `json:"bankRank,omitempty"`
	MarketItemID    *string `json:"marketItemId,omitempty"`
	MarketItemCoins *int64  `json:"marketItemCoins,omitempty"`
	IsDeciding      *bool   `json:"isDeciding,omitempty"`

	TournamentInfo *PlayerTournamentInfo `json:"tournamentInfo,omitempty"`
}

type PlayerTournamentInfo struct {
	Place       int          `json:"place"`
	Prize       int64        `json:"prize"`
	MarketPrize *MarketPrize `json:"marketPrize"`
}

func toPlayerTournamentInfo(info domain.PlayerTournamentInfo) *PlayerTournamentInfo {
	if info.Place == 0 {
		return nil
	}

	return &PlayerTournamentInfo{
		Place:       info.Place,
		Prize:       info.Prize,
		MarketPrize: toMarketPrize(info.GetTournamentMarketPrize()),
	}
}

func NewPlayer(p *domain.Player) *Player {
	var showDownAction *domain.ShowDownActionType
	if p.ShowDownAction != "" {
		showDownAction = &p.ShowDownAction
	}
	return &Player{
		UserId:          &p.UserId,
		Username:        &p.Username,
		Picture:         &p.Picture,
		Position:        &p.Position,
		Stack:           &p.Stack,
		Status:          &p.Status,
		Cards:           nil,
		Blind:           &p.Blind,
		TotalRoundBet:   &p.TotalRoundBet,
		LastGameBet:     &p.LastGameBet,
		LastGameAction:  &p.LastGameAction,
		ShowDownAction:  showDownAction,
		Level:           p.Level,
		BankBalance:     p.BankBalance,
		BankRank:        p.BankRank,
		MarketItemID:    &p.MarketItemID,
		MarketItemCoins: &p.MarketItemCoinsDayPrice,
		TournamentInfo:  toPlayerTournamentInfo(p.GetTournamentInfo()),
	}
}

func ToPlayerOpenCards(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	player := NewPlayer(p)
	player.Cards = holeCardsToCards(p.Cards)
	return player
}

func ToPlayerNoCards(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	player := NewPlayer(p)
	player.Cards = nil
	return player
}

func ToPlayerOpenLeftCard(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	player := NewPlayer(p)
	player.Cards = &[]Card{toCard(p.Cards.First), "Xx"}
	return player
}

func ToPlayerOpenRightCard(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	player := NewPlayer(p)
	player.Cards = &[]Card{"Xx", toCard(p.Cards.Second)}
	return player
}

func ToPlayerNoCardsNotDeciding(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	player := NewPlayer(p)
	player.Cards = nil
	player.IsDeciding = &False
	return player
}

func ToPlayerStackOverflow(p *domain.Player) *Player {
	if p == nil {
		return nil
	}
	player := ToPlayerNoCards(p)
	stack := p.Stack - p.TotalRoundBet
	player.Stack = &stack
	return player
}

func ToPlayerInMultiLobby(p *domain.Player) *Player {
	if p == nil {
		return nil
	}
	return &Player{
		Position: &p.Position,
		UserId:   &p.UserId,
		Username: &p.Username,
		Picture:  &p.Picture,
		Stack:    &p.StartGameStack,
	}
}

func ToPlayerOnlyPosition(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	return &Player{Position: &p.Position}
}

func ToPlayerDeciding(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	return &Player{Position: &p.Position, IsDeciding: &True}
}

func ToPlayerSecretCards(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	newP := NewPlayer(p)
	newP.Cards = &SecretHoleCards
	return newP
}

func ToPlayerGameReset(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	return &Player{
		UserId:         &p.UserId,
		Username:       &p.Username,
		Picture:        &p.Picture,
		Position:       &p.Position,
		Status:         &p.Status,
		Stack:          &p.StartGameStack,
		Blind:          &EmptyBlind,
		TotalRoundBet:  &Zero,
		LastGameBet:    &Zero,
		LastGameAction: &EmptyAction,
		Cards:          &[]Card{},
	}
}

func ToPlayerRoundReset(p *domain.Player) *Player {
	return &Player{
		Position:      &p.Position,
		TotalRoundBet: &Zero,
	}
}

func ToPlayerMultiPrevousPosition(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	player := NewPlayer(p)
	player.Cards = nil
	multiPreviousPosition := p.GetMultiPreviousPosition()
	player.Position = &multiPreviousPosition
	return player
}

func ToPlayerEnrichment(p *domain.Player) *Player {
	if p == nil {
		return nil
	}

	return &Player{
		Position:        &p.Position,
		Level:           p.Level,
		BankBalance:     p.BankBalance,
		BankRank:        p.BankRank,
		MarketItemID:    &p.MarketItemID,
		MarketItemCoins: &p.MarketItemCoinsDayPrice,
	}
}

var NillifyStack = func(p *Player) {
	p.Stack = nil
}
