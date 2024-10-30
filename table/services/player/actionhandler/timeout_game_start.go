package actionhandler

import (
	"context"
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/playerbank"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/glossd/pokergloss/table/web/client/bankclient"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func LaunchDelayedGame(staleTable *domain.Table) {
	key := timeout.Key{TableID: staleTable.ID, Position: -1, Version: staleTable.GameFlowVersion + 1}

	if conf.Props.Table.GameEndMinTimeout == 0 {
		// blocking, for tests
		DoStartGameNoCtx(key)
	} else if conf.Props.Table.GameEndMinTimeout == -1 {
		// do it yourself, for tests
	} else {
		mqpub.PublishTimeoutEvent(&timeout.Event{
			Type: timeout.StartGame,
			At:   staleTable.DecisionTimeoutAt,
			Key:  key,
		})
	}
}

// For tests
func DoStartGameNoCtx(key timeout.Key) (tryAgain bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return DoStartGame(ctx, key)
}

func DoStartGame(ctx context.Context, key timeout.Key) (tryAgain bool) {
	tableID := key.TableID
	table, err := db.FindTableGameFlow(ctx, key.TableID, key.Version)
	if err != nil {
		if err == db.ErrVersionNotMatch {
			log.Warnf("Start game, race condition, key=%s: %s", key, err)
			return false
		}
		log.Errorf("Couldn't launch next game for table, finding table, tableID=%s: %s", tableID, err)
		return true
	}

	if table.IsMultiType() {
		multiMovePlayers(table)
	}

	if table.IsCashType() {
		for _, p := range table.AutoReBuyPlayers() {
			err := bankclient.Withdraw(ctx, p.BuyInStack, p.UserId, fmt.Sprintf("Auto rebuy on table %s", table.Name))
			if err != nil {
				log.Errorf("DoStartGame: failed to rebuy: bankclient.WithdrawNoCtx: %s", err)
				continue
			}
			p.ReBuy()
		}

		for _, p := range table.AutoTopUpPlayers() {
			err := bankclient.Withdraw(ctx, p.BuyInStack-p.Stack, p.UserId, fmt.Sprintf("Auto top up on table %s", table.Name))
			if err != nil {
				log.Errorf("DoStartGame: failed to top up: bankclient.WithdrawNoCtx: %s", err)
				continue
			}
			p.TopUp()
		}
	}

	err = table.StartNextGame()
	if err != nil {
		return
	}

	err = db.SetTableGameFlow(ctx, tableID, table.GameFlowVersion, db.AllTableUpdatesGameFlow(table))
	if err != nil {
		if err == db.ErrVersionNotMatch {
			log.Errorf("Start game, set table version mismatch, tableId=%s, gameFlowVersion=%d", table.ID.Hex(), table.GameFlowVersion)
			return false
		} else {
			log.Errorf("Couldn't handle game end for table, saving table, tableID=%s: %s", tableID, err)
			return true
		}
	}

	if table.IsCashType() {
		for _, p := range table.BrokePlayers() {
			seat, err := table.GetSeat(p.Position)
			if err != nil {
				log.Errorf("Seat of player not found: failed to launch seat reservation timeout after player getting broke, tableID=%s, position=%d", table.ID.Hex(), p.Position)
				continue
			}
			LaunchSeatReservationTimeout(timeout.Key{TableID: table.ID, Position: seat.Position, Version: seat.Version}, seat.ReservationTimeoutAt())
		}
	}

	if table.IsWaiting() {
		broadcast.SendTableEvents(tableID.Hex(), append(playerbank.HandleNullifiedPlayersLeft(table), events.BuildReset(table)))
		switch table.Type {
		case domain.SitngoType:
			_ = db.UpdateSitngoLobbyStatus(table.LobbyID, domain.LobbyFinished)
			_ = db.DeleteTableNoCtx(table.ID)
		case domain.MultiType:
			broadcast.SendMultiTablePlayersUpdates(table)
			if table.IsLast {
				_ = db.UpdateLobbyMultiStatusNoCtx(table.LobbyID, domain.LobbyFinished)
				_ = db.DeleteTableNoCtx(table.ID)
			}
		}
		if table.IsSurvival {
			mqpub.PublishSurvivalEnd(table)
		}
	} else {
		HandleGameStart(table)
		if table.IsMultiType() {
			broadcast.SendMultiTablePlayersUpdates(table)
		}
	}
	return
}

func multiMovePlayers(tableFrom *domain.Table) {
	toTablesMap, err := fetchMoveToTables(tableFrom)

	if err != nil || len(toTablesMap) == 0 {
		return
	}

	for _, move := range tableFrom.PlayerMoves {
		s, err := tableFrom.GetSeat(move.FromPosition)
		if err == nil && s.IsTaken() {
			p := s.GetPlayer()
			t, ok := toTablesMap[move.ToTableID]
			if ok {
				err := t.MultiPutPlayerAtFreePosition(p)
				if err != nil {
					// todo maybe move him to another table
					log.Errorf("Couldn't move player: %s", err)
					continue
				}
				tableFrom.MultiDeletePlayerForMove(s, t.ID)
			} else {
				log.Errorf("Table not found to move player to another table, toTableId=%s, tableMap=%+v", move.ToTableID.Hex(), toTablesMap)
			}
		}
	}

	toTables := make([]*domain.Table, 0, len(toTablesMap))
	for _, t := range toTablesMap {
		err = db.SetTableMultiPutPlayers(t)
		if err != nil {
			log.Errorf("Failed to save put players on table: %s", err)
			continue
		}
		toTables = append(toTables, t)
	}

	tableFrom.PlayerMoves = nil
	_ = db.NullifyTablePlayerMoves(tableFrom.ID)

	broadcast.SendEventsAboutMovingPlayers(tableFrom, toTables)
	mqpub.PublishMultiRebalance(tableFrom.LobbyID.Hex(), tableFrom.ID.Hex())
}

func fetchMoveToTables(table *domain.Table) (map[primitive.ObjectID]*domain.Table, error) {
	toTableIDsMap := make(map[primitive.ObjectID]struct{})
	for _, move := range table.PlayerMoves {
		toTableIDsMap[move.ToTableID] = struct{}{}
	}
	uniqueTableToIDs := make([]primitive.ObjectID, 0, len(toTableIDsMap))
	for id := range toTableIDsMap {
		uniqueTableToIDs = append(uniqueTableToIDs, id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	tables, err := db.FilterTables(ctx, bson.M{"_id": bson.M{"$in": uniqueTableToIDs}})
	if err != nil {
		log.Errorf("Failed move players, find failed: %s", err)
		return nil, err
	}
	tablesMap := make(map[primitive.ObjectID]*domain.Table)
	for _, table := range tables {
		tablesMap[table.ID] = table
	}
	return tablesMap, nil
}

func HandleGameStart(table *domain.Table) {
	var beforeEventsAcc []*events.TableEvent
	beforeEventsAcc = append(beforeEventsAcc, events.BuildReset(table))
	beforeEventsAcc = append(beforeEventsAcc, playerbank.HandleNullifiedPlayersLeft(table)...)
	beforeEventsAcc = append(beforeEventsAcc, events.BuildBlinds(table))

	userEvents := events.BuildUserHoleCards(table)
	notFoundUserEvents := []*events.TableEvent{events.BuildTableHoleCardsAllSecret(table)}
	secretEvents := []*events.TableEvent{events.BuildAllFaceUpHoleCards(table)}

	var afterEventsAcc []*events.TableEvent
	// case forced blinds
	if table.IsGameEnd() {
		afterEventsAcc = append(afterEventsAcc, events.BuildNewBettingRound(table))
		afterEventsAcc = append(afterEventsAcc, GameEnd{}.WsEvents(table)...)
		GameEnd{}.Timeout(table)
	} else {
		LaunchDecisionTimeout(table)
		afterEventsAcc = append(afterEventsAcc, events.BuildTimeToDecide(table))
	}

	mqpub.SendTableMessageToUsers(table.ID.Hex(), userEvents, notFoundUserEvents, secretEvents, beforeEventsAcc, afterEventsAcc)
}
