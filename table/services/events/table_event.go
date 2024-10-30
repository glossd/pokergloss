package events

import (
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
)

// It's 'tableless', event doesn't belong to any specific table
type TableEvent struct {
	Type TET `json:"type" enums:"initState,seatReserved,seatReservationTimeout,bankroll,playerLeft,blinds,holeCards,timeToDecide,timeToDecideTimeout,playerAction,newBettingRound,showDown,winners,reset,sitBack,addChips,intent,playerMoved,enrichPlayers,stackOverflowPlayer,setPlayerConfig,multiPlayersUpdate,multiPlusPlayersUpdate"`
	// Exclusively for skipping an event for the user himself/herself
	Payload interface{} `json:"payload"`
}

// Shortcut, map[string]interface{} is too long to write
type M map[string]interface{}

// TableEvent Type
type TET string

// ATTENTION!!! when changing list of TETs don't forget to edit enums of TableEvent.Type
const (
	// The state of the table, emits on the start of user connection
	InitState TET = "initState"
	// A seat is reserved by a player
	SeatReserved           TET = "seatReserved"
	SeatReservationTimeout TET = "seatReservationTimeout"
	// A player with reserved seat put his chips in the game and started playing
	Bankroll   TET = "bankroll"
	PlayerLeft TET = "playerLeft"
	// Tells who is bigBlind, smallBlind and dealer
	Blinds TET = "blinds"
	// Deals hole cards to table players
	HoleCards TET = "holeCards"
	// It's time to decide for player which betting action to make
	TimeToDecide        TET = "timeToDecide"
	TimeToDecideTimeout TET = "timeToDecideTimeout"
	// Tells what kind of betting action player has made
	PlayerMadeAction TET = "playerAction"
	NewBettingRound  TET = "newBettingRound"
	ShowDown         TET = "showDown"
	Winners          TET = "winners"
	Reset            TET = "reset"

	PlayerSitBack TET = "sitBack"
	AddChips      TET = "addChips"
	Intent        TET = "intent"

	SetPlayerConfig TET = "setPlayerConfig"

	StackOverflowPlayer TET = "stackOverflowPlayer"
)

// TableEvent Payload
type TEP struct {
	Table *model.Table `json:"table"`
}

func BuildInitState(table *domain.Table, userID string) *TableEvent {
	return &TableEvent{Type: InitState, Payload: TEP{Table: model.ToModelTableSeats(table, model.AllSeatsIdentifiedCards(table, userID))}}
}

func BuildSeatReservedEvent(seat *domain.Seat) *TableEvent {
	return &TableEvent{Type: SeatReserved, Payload: TEP{Table: model.TableSeat(model.ToSeat(seat, model.ToPlayerNoCards))}}
}

func BuildSeatReservationTimeout(position int) *TableEvent {
	return &TableEvent{Type: SeatReservationTimeout, Payload: TEP{Table: model.TableSeat(model.EmptySeat(position))}}
}

func BuildBankroll(player *domain.Player) *TableEvent {
	return &TableEvent{Type: Bankroll, Payload: TEP{Table: model.TableSeat(model.PlayerToSeat(player, model.ToPlayerNoCards))}}
}

func BuildPlayerLeft(player *domain.Player) *TableEvent {
	return &TableEvent{Type: PlayerLeft, Payload: M{
		"table":      model.TableSeat(model.EmptySeat(player.Position)),
		"leftPlayer": model.ToPlayerNoCards(player),
	}}
}

func BuildBlinds(table *domain.Table) *TableEvent {
	var seats []*model.Seat
	for _, p := range table.UniqueBlindsSlice() {
		if p == nil {
			continue
		}
		seats = append(seats, model.PlayerToSeat(p, func(p *domain.Player) *model.Player {
			res := model.ToPlayerNoCards(p)
			if table.IsAutoGameEnd() {
				stack := p.StartGameStack - p.LastGameBet
				res.Stack = &stack
			}
			return res
		}))
	}
	return &TableEvent{Type: Blinds, Payload: TEP{Table: model.TableBlinds(table, seats)}}
}

func BuildTableHoleCards(table *domain.Table, userId string) *TableEvent {
	return &TableEvent{Type: HoleCards, Payload: TEP{
		Table: &model.Table{Seats: SeatsHoleCards(table, userId)},
	}}
}

func BuildTableHoleCardsAllSecret(table *domain.Table) *TableEvent {
	return &TableEvent{Type: HoleCards, Payload: TEP{
		Table: &model.Table{Seats: model.PlayersToSeats(table.PlayingPlayersByGameType(), model.ToPlayerSecretCards, model.NillifyStack)},
	}}
}

func BuildAllFaceUpHoleCards(table *domain.Table) *TableEvent {
	return &TableEvent{Type: HoleCards, Payload: TEP{
		Table: &model.Table{Seats: model.PlayersToSeats(table.PlayingPlayersByGameType(), model.ToPlayerOpenCards, model.NillifyStack)},
	}}
}

func BuildTimeToDecide(table *domain.Table) *TableEvent {
	return TimeToDecideBuilder(table, false)
}

func TimeToDecideBuilder(table *domain.Table, withTableStatus bool) *TableEvent {
	var tableStatus *domain.TableStatus
	if withTableStatus {
		tableStatus = &table.Status
	}

	decidingPlayer, err := table.DecidingPlayer()
	if err != nil {
		return nil
	}
	return &TableEvent{Type: TimeToDecide, Payload: M{
		"table": model.Table{
			Seats:                 []*model.Seat{model.PlayerToSeatTimeout(decidingPlayer, table.DecisionTimeoutAt, model.ToPlayerDeciding)},
			DecidingPosition:      &table.DecidingPosition,
			LastAggressorPosition: &table.LastAggressorPosition,
			Status:                tableStatus,
		},
	}}
}

func BuildNewBettingRound(table *domain.Table) *TableEvent {
	return &TableEvent{Type: NewBettingRound, Payload: M{
		"table": model.Table{
			TotalPot:       &table.TotalPot,
			Pots:           model.ToPots(table.Pots),
			CommunityCards: model.ToCards(table.CommunityCards.AvailableCards()),
			Seats:          model.PlayersToSeats(table.AllPlayers(), model.ToPlayerRoundReset),
		},
		"roundType": table.CommunityCards.RoundType(),
		"newCards":  model.ToCards(table.CommunityCards.GetNewCards()),
	}}
}

func BuildWinners(table *domain.Table) *TableEvent {
	return &TableEvent{Type: Winners, Payload: M{
		"table": model.Table{
			Status:           &table.Status,
			Pots:             model.ToPots(table.Pots),
			TotalPot:         &table.TotalPot,
			DecidingPosition: &model.NegativeOne,
			CommunityCards:   model.ToCards(table.CommunityCards.AvailableCards()),
			Seats:            model.PlayersToSeats(table.AllPlayers(), model.ToPlayerRoundReset),
			Winners:          model.ToWinners(table),
			Rakes:            model.RakeToUserRakes(table.GetRake()),
		},
		"newCards": model.ToCards(table.CommunityCards.GetNewCards()),
	}}
}

func BuildReset(table *domain.Table) *TableEvent {
	return &TableEvent{Type: Reset, Payload: TEP{Table: model.TableReset(table)}}
}

func BuildPlayerSitBack(player *domain.Player) *TableEvent {
	return &TableEvent{Type: PlayerSitBack,
		Payload: TEP{Table: model.TableSeat(model.PlayerToSeat(player, func(p *domain.Player) *model.Player {
			return &model.Player{Position: &p.Position, Status: &p.Status}
		}))},
	}
}

func BuildAddChips(player *domain.Player) *TableEvent {
	return &TableEvent{Type: AddChips, Payload: TEP{
		Table: model.TableSeat(model.PlayerToSeat(player, func(p *domain.Player) *model.Player {
			return &model.Player{Position: &p.Position, Stack: &p.Stack}
		}))},
	}
}

func BuildIntent(player *domain.Player) *TableEvent {
	return &TableEvent{Type: Intent, Payload: TEP{
		Table: model.TableSeat(model.PlayerToSeat(player, func(p *domain.Player) *model.Player {
			return &model.Player{Position: &p.Position, Intent: model.ToIntent(player.Intent)}
		})),
	}}
}

func BuildShowDown(player *domain.Player, decPos int) *TableEvent {
	var mapper model.PlayerMapper
	switch player.ShowDownAction {
	case domain.Muck:
		mapper = model.ToPlayerNoCards
	case domain.ShowRight:
		mapper = model.ToPlayerOpenRightCard
	case domain.ShowLeft:
		mapper = model.ToPlayerOpenLeftCard
	case domain.Show:
		mapper = model.ToPlayerOpenCards
	}
	table := model.TableSeat(model.PlayerToSeat(player, mapper, func(p *model.Player) {
		p.Stack = nil
		p.IsDeciding = &model.False
	}))
	table.DecidingPosition = &decPos
	return &TableEvent{Type: ShowDown, Payload: TEP{Table: table}}
}

func BuildStackOverflowPlayer(player *domain.Player) *TableEvent {
	return &TableEvent{
		Type:    StackOverflowPlayer,
		Payload: TEP{Table: model.TableSeat(model.PlayerToSeat(player, model.ToPlayerNoCards))},
	}
}

func BuildSetPlayerConfig(config *domain.PlayerAutoConfig) []*TableEvent {
	evts := []*TableEvent{{
		Type:    SetPlayerConfig,
		Payload: M{"config": model.ToPlayerAutoConfig(config)},
	}}
	return evts
}
