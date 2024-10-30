package db

import (
	"context"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertTable(ctx context.Context, table *domain.Table) error {
	_, err := ColTable().InsertOne(ctx, table)
	if err != nil {
		log.Errorf("Insert table failed: %s", err)
		return err
	}
	return nil
}

func SaveTableNoCtx(table *domain.Table) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return InsertTable(ctx, table)
}

func InsertManyTables(ctx context.Context, tables []*domain.Table) error {
	adaptedTables := make([]interface{}, 0, len(tables))
	for _, table := range tables {
		adaptedTables = append(adaptedTables, table)
	}

	_, err := ColTable().InsertMany(ctx, adaptedTables)
	if err != nil {
		log.Errorf("Failed insert many tables: %s", err)
		return err
	}
	return nil
}

func UpdateTableContext(ctx context.Context, table *domain.Table) error {
	_, err := ColTable().ReplaceOne(ctx, filterID(table.ID), table)
	if err != nil {
		log.Errorf("Couldn't update table, tableID=%s : %s", table.ID, err)
		return err
	}
	return nil
}

func DeleteTable(ctx context.Context, tableID primitive.ObjectID) error {
	_, err := ColTable().DeleteOne(ctx, bson.D{{"_id", tableID}})
	if err != nil {
		log.Errorf("Delete table %s failed: %s", tableID.Hex(), err)
		return err
	}
	return nil
}

func DeleteTableNoCtx(tableID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return DeleteTable(ctx, tableID)
}

func SetTableReservation(ctx context.Context, table *domain.Table, seat *domain.Seat) error {
	filter := bson.D{{"_id", table.ID}, SeatVersion(seat)}
	seat.Version++
	update := bson.D{{"$set",
		bson.D{{SeatDbPath(seat.Position), seat}, PlayersCount(table)},
	}}
	one, err := ColTable().UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("Update table %s reservation failed: %s", table.ID.Hex(), err)
		return err
	}
	if one.MatchedCount == 0 {
		log.Warnf("Seat reservation, race condition happened, tableID=%s, position=%d", table.ID.Hex(), seat.Position)
		return ErrVersionNotMatch
	}

	return nil
}

func CancelTableReservation(ctx context.Context, table *domain.Table, seat *domain.Seat) error {
	tableID := table.ID
	filter := bson.D{{"_id", tableID}, SeatVersion(seat)}
	seat.Version++
	update := bson.D{{"$set", bson.D{PlayerNullify(seat.Position), IncSeatVersion(seat), PlayersCount(table)}}}
	one, err := ColTable().UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("Failed to set table %s, cancel reservation: %s", tableID.Hex(), err)
		return err
	}
	if one.MatchedCount == 0 {
		log.Warnf("Seat reservation, race condition happened, tableID=%s, position=%d", tableID.Hex(), seat.Position)
		return ErrVersionNotMatch
	}
	return nil
}

func SetTableMultiAttrs(table *domain.Table) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	update := bson.D{{"$set", bson.D{{"multiattrs", table.MultiAttrs}}}}
	_, err := ColTable().UpdateOne(ctx, filterID(table.ID), update)
	if err != nil {
		log.Errorf("Failed to update table multi attrs, tableId=%s: %s", table.ID.Hex(), err)
		return err
	}
	return nil
}

func NullifyTablePlayerMoves(tableID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	update := bson.D{{"$set", bson.D{{"multiattrs.playermoves", nil}}}}
	_, err := ColTable().UpdateOne(ctx, filterID(tableID), update)
	if err != nil {
		log.Errorf("Failed to nullify table PlayerMoves, tableId=%s: %s", tableID.Hex(), err)
		return err
	}
	return nil
}

func SetTableMultiPutPlayers(table *domain.Table) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	var update []bson.E
	for _, player := range table.MultiAttrs.GetPutPlayers() {
		update = append(update, bson.E{Key: PlayerDbPath(player.Position), Value: player})
	}

	_, err := ColTable().UpdateOne(ctx, filterID(table.ID), bson.M{"$set": update})
	if err != nil {
		log.Errorf("Failed to update table multi put players, tableId=%s: %s", table.ID.Hex(), err)
		return err
	}
	return nil
}
