package events

import (
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	MultiPlayerMoveCommandType TET = "multiPlayerMove"
	MultiPlayersUpdateType     TET = "multiPlayersUpdate"
	MultiPlusPlayersUpdateType TET = "multiPlusPlayersUpdate"
	PlayerMovedType            TET = "playerMoved"
)

func BuildMultiPlayersUpdate(table *domain.Table) *TableEvent {
	return &TableEvent{Type: MultiPlayersUpdateType, Payload: M{
		"tableId": table.ID.Hex(),
		"players": model.ToUsersWithStack(table.AllPlayers()),
	}}
}

func BuildMultiPlusPlayersUpdate(tableID primitive.ObjectID, players []*domain.Player) *TableEvent {
	return &TableEvent{Type: MultiPlusPlayersUpdateType, Payload: M{
		"tableId": tableID.Hex(),
		"players": model.ToUsersWithStack(players),
	}}
}

func BuildMultiPlayersEmpty(tableID primitive.ObjectID) *TableEvent {
	return &TableEvent{Type: MultiPlayersUpdateType, Payload: M{
		"tableId": tableID.Hex(),
		"players": []model.UserWithStack{},
	}}
}

func BuildMultiPlayerMove(tableIdToMoveTo primitive.ObjectID) *TableEvent {
	return &TableEvent{
		Type:    MultiPlayerMoveCommandType,
		Payload: M{"tableId": tableIdToMoveTo},
	}
}

func BuildPlayerMoved(player *domain.Player) *TableEvent {
	return &TableEvent{Type: PlayerMovedType, Payload: M{
		"table":      model.TableSeat(model.EmptySeat(player.GetMultiPreviousPosition())),
		"leftPlayer": model.ToPlayerMultiPrevousPosition(player),
	}}
}
