package player

import (
	"errors"
	"fmt"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/enrich"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/web/client/bankclient"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrBuyInRaceCondition = services.ErrFormat("state of your table has changed too fast, please buy in again")

func BuyIn(params *ChipsParams) (bool, error) {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		log.Errorf("Couldn't buy in for %s: %s", params.iden, err)
		return false, err
	}

	return BuyInOnTable(table, params)
}

func BuyInOnTable(table *domain.Table, params *ChipsParams) (bool, error) {
	ctx := params.ctx
	tableID := params.tableID
	stack := params.chips
	position := params.position
	iden := params.iden

	isNewGame, err := table.BuyIn(stack, position, iden)
	if err != nil {
		return false, err
	}

	buyInSeat := table.GetSeatUnsafe(position)

	err = bankclient.Withdraw(ctx, stack, iden.UserId, fmt.Sprintf("Joined table %s", table.Name))
	if err != nil {
		if err == bankclient.ErrNotEnoughChips {
			return false, err
		}
		log.Errorf("Couldn't buy in for %s, bank service error: %s", iden, err)
		return false, fmt.Errorf("bank service unavailable: %s", err)
	}

	var dbUpdates []primitive.E
	var dbErr error
	if isNewGame || table.IsWaiting() { // mongo can't update the same elements
		dbUpdates = db.AllTableUpdatesGameFlow(table)
		dbUpdates = append(dbUpdates, db.IncSeatVersion(buyInSeat))
		filter := []bson.E{db.GameFlowVersion(table), db.SeatVersion(buyInSeat)}
		dbErr = db.SetTableFilter(ctx, table.ID, filter, dbUpdates)
	} else {
		dbUpdates = db.SeatUpdate(buyInSeat)
		dbUpdates = append(dbUpdates, db.IncSeatVersion(buyInSeat))
		dbErr = db.SetTableContext(ctx, table.ID, dbUpdates)
	}

	if dbErr != nil {
		bankclient.Deposit(stack, iden.UserId, fmt.Sprintf("Failed to join table %s", table.Name))
		if errors.Is(dbErr, mongo.ErrNoDocuments) {
			return false, ErrBuyInRaceCondition
		} else {
			log.Errorf("Couldn't buy in for %s: %s", iden, dbErr)
		}
		return false, err
	}

	broadcast.SendTableEvent(tableID.Hex(), events.BuildBankroll(table.GetPlayerUnsafe(position)))

	if isNewGame {
		actionhandler.HandleGameStart(table)
	}

	go enrich.Players(table, []*domain.Player{buyInSeat.GetPlayer()})

	return isNewGame, nil
}
