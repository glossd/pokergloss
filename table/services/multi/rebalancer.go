package multi

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
)

type RebalanceResult struct {
	Status                       RebalanceStatus
	CountTablesWithMovingPlayers int
}

type RebalanceStatus string

const (
	StopRebalance   RebalanceStatus = "stop"
	RemoveRightAway RebalanceStatus = "removed"
	MoveAllPlayers  RebalanceStatus = "moveAllPlayers"
	Disproportion   RebalanceStatus = "disproportion"
)

func Rebalance(lobbyID primitive.ObjectID) (*RebalanceResult, error) {
	tables, err := db.FindTablesByLobbyID(lobbyID)
	if err != nil {
		log.Errorf("Failed rebalance players, find tables lobbyID=%s, %s", lobbyID.Hex(), err)
		return nil, err
	}

	if len(tables) < 2 {
		return &RebalanceResult{Status: StopRebalance}, nil
	}

	for _, t := range tables {
		t.MultiAttrs.PlayerMoves = nil
	}

	var sumPlayers int
	for _, table := range tables {
		table.PlayersCount = len(table.AllPlayers())
		sumPlayers += table.PlayersCount
	}

	// From most full
	sort.Slice(tables, func(i, j int) bool {
		return tables[i].PlayersCount > tables[j].PlayersCount
	})

	var result RebalanceResult
	if isFewEnoughToRemoveTable(tables, sumPlayers) {
		result = rebalanceByRemovingTable(tables)
	} else {
		rebalanceDisproportionalTables(tables, sumPlayers)

		var count int
		for _, table := range tables {
			if len(table.PlayerMoves) > 0 {
				count++
			}
		}
		result = RebalanceResult{
			Status:                       Disproportion,
			CountTablesWithMovingPlayers: count,
		}
	}

	for _, table := range tables {
		_ = db.SetTableMultiAttrs(table)
	}

	return &result, nil
}

func isFewEnoughToRemoveTable(sortedTables []*domain.Table, sumPlayers int) bool {
	tableSize := sortedTables[0].Size
	tablesNum := len(sortedTables)
	return tableSize*(tablesNum-1) >= sumPlayers
}

func rebalanceDisproportionalTables(tables []*domain.Table, sumPlayers int) {
	// moves players on game end
	avgPlayersCount := float64(sumPlayers) / float64(len(tables))
	for i := len(tables) - 1; i >= 0; i-- { // reverse for loop
		if needsMorePlayers(tables, i) {
			tableToMoveTo := tables[i]
			howManyPlayersToMove := int(avgPlayersCount) - tableToMoveTo.PlayersCount
			for j := 0; j < howManyPlayersToMove; j++ {
				addPlayerToGameEndMove(tables[:i], tableToMoveTo, avgPlayersCount)
			}
		} else {
			break
		}
	}
}

func needsMorePlayers(tables []*domain.Table, idx int) bool {
	var sumPlayers int
	for _, table := range tables {
		sumPlayers += table.PlayersCount
	}
	avgPlayersCount := float64(sumPlayers) / float64(len(tables))
	if idx < 0 || idx >= len(tables) {
		log.Errorf("needsMorePlayers index out of boundaries, tablesLen=%d, idx=%d", len(tables), idx)
		return false
	}
	return float64(tables[idx].PlayersCount+1) <= avgPlayersCount
}

func addPlayerToGameEndMove(tablesToRemoveFrom []*domain.Table, tableMoveTo *domain.Table, avgPlayerCount float64) {
	if len(tablesToRemoveFrom) == 0 {
		return
	}
	tableToRemoveFrom := tablesToRemoveFrom[0]
	if float64(tableToRemoveFrom.PlayersCount) >= avgPlayerCount+1 {
		tableToRemoveFrom.MultiAddPlayerToMove(tableToRemoveFrom.MultiRandomAvailablePlayerToMove(), tableMoveTo)
		return
	} else {
		if len(tablesToRemoveFrom) > 1 {
			addPlayerToGameEndMove(tablesToRemoveFrom[1:], tableMoveTo, avgPlayerCount)
		}
	}
}

func rebalanceByRemovingTable(tables []*domain.Table) RebalanceResult {
	tableToRemove := tables[len(tables)-1]
	tableToRemove.PlayerMoves = nil

	tablesToKeep := tables[:len(tables)-1]

	if tableToRemove.IsWaiting() {
		removeRightAway(tableToRemove, tablesToKeep)
		return RebalanceResult{Status: RemoveRightAway, CountTablesWithMovingPlayers: 0}
	} else {
		moveAllPlayerFrom(tableToRemove, tablesToKeep)
		return RebalanceResult{Status: MoveAllPlayers, CountTablesWithMovingPlayers: 1}
	}
}

// Moves players on game end.
func moveAllPlayerFrom(fromTable *domain.Table, tablesToKeep []*domain.Table) {
	fromPlayers := fromTable.AllPlayers()
	for _, table := range tablesToKeep {
		if len(fromPlayers) == 0 {
			break
		}
		availableSeatsCount := table.Size - table.PlayersCount
		for i := 0; i < availableSeatsCount; i++ {
			if len(fromPlayers) > 0 {
				playerToMove := fromPlayers[len(fromPlayers)-1]
				fromTable.MultiAddPlayerToMove(playerToMove, table)
				fromPlayers = fromPlayers[:len(fromPlayers)-1]
			}
		}
	}
}

func removeRightAway(tableToRemove *domain.Table, tablesToKeep []*domain.Table) {
	// todo there's maybe multiple tables with
	if tableToRemove.PlayersCount == 0 {
		err := db.DeleteTableNoCtx(tableToRemove.ID)
		if err != nil {
			log.Errorf("Couldn't delete table of multi lobby, tableId=%s : %s", tableToRemove.ID, err)
		}
		checkAndUpdateLast(tablesToKeep)
		return
	}
	// Deletes table right away
	p, onlyOne := tableToRemove.IsOnlyOneOnTable()
	if !onlyOne {
		log.Errorf("Table is waiting with multiple players tableId=%s", tableToRemove.ID.Hex())
		return
	}

	tableToUpdate := tablesToKeep[len(tablesToKeep)-1]
	isNewGame, err := tableToUpdate.MultiSitPlayerAndTryToStartGame(p)
	if err != nil {
		log.Errorf("Failed to rebalance: %s", err)
		return
	}

	checkAndUpdateLast(tablesToKeep)

	err = db.DeleteTableNoCtx(tableToRemove.ID)
	if err != nil {
		log.Errorf("Couldn't delete table of multi lobby, tableId=%s : %s", tableToRemove.ID, err)
	}

	if isNewGame {
		err := db.SetTableGameFlowNoCtx(tableToUpdate.ID, tableToUpdate.GameFlowVersion, db.AllTableUpdatesGameFlow(tableToUpdate))
		if err != nil {
			log.Errorf("Failed to rebalance: %s", err)
			return
		}
	} else {
		err = db.SetTableMultiPutPlayers(tableToUpdate)
		if err != nil {
			log.Errorf("Failed to rebalance: %s", err)
			return
		}
	}

	broadcast.SendEventsAboutMovingPlayers(tableToRemove, []*domain.Table{tableToUpdate})

	// redirects everyone
	broadcast.SendTableEvent(tableToRemove.ID.Hex(), events.BuildMultiPlayerMove(tableToUpdate.ID))

	if isNewGame {
		actionhandler.LaunchDelayedGame(tableToUpdate)
	}

	tableEvents := []*events.TableEvent{
		events.BuildMultiPlayersEmpty(tableToRemove.ID),
		events.BuildMultiPlayersUpdate(tableToUpdate),
	}

	for _, table := range tablesToKeep {
		broadcast.SendTableEvents(table.ID.Hex(), tableEvents)
	}
}

func checkAndUpdateLast(tablesToKeep []*domain.Table) {
	if len(tablesToKeep) == 1 {
		tablesToKeep[0].MultiAttrs.IsLast = true
	}
}
