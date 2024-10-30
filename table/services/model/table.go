package model

import (
	"github.com/glossd/pokergloss/table/domain"
)

type Table struct {
	ID                 *string              `json:"id,omitempty"`
	Name               *string              `json:"name,omitempty"`
	Size               *int                 `json:"size,omitempty"`
	BigBlind           *int64               `json:"bigBlind,omitempty"`
	SmallBlind         *int64               `json:"smallBlind,omitempty"`
	DecisionTimeoutSec *int                 `json:"decisionTimeoutSec,omitempty"`
	BettingLimit       *domain.BettingLimit `json:"bettingLimit,omitempty" enums:"NL,PL,FL,ML"`
	MinBuyIn           *int64               `json:"minBuyIn,omitempty"`
	MaxBuyIn           *int64               `json:"maxBuyIn,omitempty"`

	Type *domain.TableType `json:"type,omitempty" enums:"cashGame,sitAndGo,multi"`

	Occupied *int   `json:"occupied,omitempty"`
	AvgStake *int64 `json:"avgStake,omitempty"`
	AvgPot   *int64 `json:"avgPot,omitempty"`

	Seats []*Seat `json:"seats,omitempty"`

	Status *domain.TableStatus `json:"status,omitempty" enums:"waiting,playing,gameEnd,showDown"`

	MaxRoundBet       *int64 `json:"maxRoundBet,omitempty"`
	BettingLimitChips *int64 `json:"bettingLimitChips,omitempty"`

	Pots           *[]Pot      `json:"pots,omitempty"`
	TotalPot       *int64      `json:"totalPot,omitempty"`
	CommunityCards *[]Card     `json:"communityCards,omitempty"`
	Rakes          *[]UserRake `json:"rakes,omitempty"`

	Winners *[]*Winner `json:"winners,omitempty"`

	DecidingPosition      *int `json:"decidingPosition,omitempty"`
	LastAggressorPosition *int `json:"lastAggressorPosition,omitempty"`

	TournamentAttrs *TournamentAttrs `json:"tournamentAttrs,omitempty"`
	MultiAttrs      *MultiAttrs      `json:"multiAttrs,omitempty"`

	ThemeID       *domain.ThemeID `json:"themeId,omitempty"`
	IsSurvival    *bool           `json:"isSurvival,omitempty"`
	SurvivalLevel *int64          `json:"survivalLevel,omitempty"`
}

type Pot struct {
	Chips           int64 `json:"chips"`
	WinnerPositions []int `json:"winnerPositions"`
}

type UserRake struct {
	Position int   `json:"position"`
	Chips    int64 `json:"chips"`
}

func ToModelTables(tables []*domain.Table) []*Table {
	models := make([]*Table, len(tables))
	for i, table := range tables {
		models[i] = ToModelTable(table, ToPlayerInMultiLobby)
	}
	return models
}

func ToModelTable(t *domain.Table, mapper PlayerMapper) *Table {
	return ToModelTableSeats(t, ToSeats(t.Seats, mapper))
}

func ToModelTableSeats(t *domain.Table, seats []*Seat) *Table {
	id := t.ID.Hex()
	occupied := len(t.AllPlayers())
	avgStake := int64(0)
	avgPot := int64(0)

	minStack := t.MinBuyInStack()
	maxStack := t.MaxBuyInStack()
	maxRoundBet := t.MaxRoundBet()
	blc := t.BettingLimitChips()
	decTimout := int(t.DecisionTimeout.Seconds())
	return &Table{
		ID:                    &id,
		Name:                  &t.Name,
		Size:                  &t.Size,
		Type:                  &t.Type,
		Occupied:              &occupied,
		BigBlind:              &t.BigBlind,
		SmallBlind:            &t.SmallBlind,
		BettingLimit:          &t.BettingLimit,
		AvgStake:              &avgStake,
		AvgPot:                &avgPot,
		Seats:                 seats,
		Status:                &t.Status,
		DecisionTimeoutSec:    &decTimout,
		MinBuyIn:              &minStack,
		MaxBuyIn:              &maxStack,
		MaxRoundBet:           &maxRoundBet,
		BettingLimitChips:     &blc,
		Pots:                  ToPots(t.Pots),
		TotalPot:              &t.TotalPot,
		CommunityCards:        ToCards(t.CommunityCards.AvailableCards()),
		DecidingPosition:      &t.DecidingPosition,
		LastAggressorPosition: &t.LastAggressorPosition,
		TournamentAttrs:       toTournamentAttrs(t.TournamentAttributes),
		ThemeID:               &t.ThemeID,
		IsSurvival:            &t.IsSurvival,
		SurvivalLevel:         &t.SurvivalLevel,
	}
}

func ToPots(pots domain.Pots) *[]Pot {
	if len(pots) == 0 {
		return nil
	}

	var result []Pot
	for _, pot := range pots {
		if pot != nil && pot.Chips != 0 {
			result = append(result, Pot{Chips: pot.Chips, WinnerPositions: pot.WinnerPositions})
		}
	}
	return &result
}

func TableSeat(s *Seat) *Table {
	return &Table{Seats: []*Seat{s}}
}

func TablePlayers(players []*domain.Player, mapper PlayerMapper) *Table {
	return &Table{Seats: PlayersToSeats(players, mapper)}
}

func TableReset(table *domain.Table) *Table {
	seats := ToSeats(table.Seats, ToPlayerGameReset)
	return &Table{
		Seats:                 seats,
		Status:                &table.Status,
		Pots:                  &[]Pot{},
		Rakes:                 &[]UserRake{},
		MaxRoundBet:           &Zero,
		BettingLimitChips:     &Zero,
		CommunityCards:        &[]Card{},
		Winners:               &[]*Winner{},
		LastAggressorPosition: &NegativeOne,
		DecidingPosition:      &NegativeOne,
		TournamentAttrs:       toTournamentAttrs(table.TournamentAttributes),
	}
}

func TableBlinds(t *domain.Table, seats []*Seat) *Table {
	maxRoundBet := t.MaxRoundBet()
	blc := t.BettingLimitChips()
	return &Table{
		BigBlind:          &t.BigBlind,
		SmallBlind:        &t.SmallBlind,
		Seats:             seats,
		Status:            &t.Status,
		MaxRoundBet:       &maxRoundBet,
		BettingLimitChips: &blc,
		TotalPot:          &t.TotalPot,
		TournamentAttrs:   toTournamentAttrs(t.TournamentAttributes),
	}
}

func AllSeatsIdentifiedCards(table *domain.Table, userID string) []*Seat {
	var seats []*Seat
	for _, seat := range table.Seats {
		if seat.IsFree() {
			seats = append(seats, EmptySeat(seat.Position))
		} else {
			player := seat.GetPlayer()
			var modelSeat *Seat
			if player.HasCards() {
				if player.UserId == userID {
					modelSeat = PlayerToSeat(player, ToPlayerOpenCards)
				} else {
					modelSeat = PlayerToSeat(player, ToPlayerSecretCards)
				}
			} else {
				modelSeat = PlayerToSeat(player, ToPlayerNoCards)
			}
			if table.IsDeciding(player) {
				modelSeat.Player.TimeoutAt = &table.DecisionTimeoutAt
				modelSeat.Player.IsDeciding = &True
			}
			seats = append(seats, modelSeat)
		}
	}
	return seats
}

func AllSeatsFaceDownCards(table *domain.Table, userID string) []*Seat {
	var seats []*Seat
	for _, seat := range table.Seats {
		if seat.IsFree() {
			seats = append(seats, EmptySeat(seat.Position))
		} else {
			player := seat.GetPlayer()
			var modelSeat *Seat
			if player.HasCards() {
				if player.UserId == userID {
					modelSeat = PlayerToSeat(player, ToPlayerOpenCards)
				} else {
					modelSeat = PlayerToSeat(player, ToPlayerSecretCards)
				}
			} else {
				modelSeat = PlayerToSeat(player, ToPlayerNoCards)
			}
			if table.IsDeciding(player) {
				modelSeat.Player.TimeoutAt = &table.DecisionTimeoutAt
			}
			seats = append(seats, modelSeat)
		}
	}
	return seats
}

func RakeToUserRakes(rake domain.Rake) *[]UserRake {
	var result = make([]UserRake, 0, len(rake.PositionChips))
	for pos, chips := range rake.PositionChips {
		result = append(result, UserRake{Position: pos, Chips: chips})
	}
	return &result
}
