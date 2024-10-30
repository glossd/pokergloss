package db

import (
	"context"
	"github.com/glossd/pokergloss/bonus/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func ResetBonuses() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := BonusCol().UpdateMany(ctx, bson.M{"istaken": false}, bson.M{"$set": bson.M{"dayinarow": 0}})
	if err != nil {
		log.Errorf("Couldn't update bonuses: %s", err)
		return
	}

	_, err = BonusCol().UpdateMany(ctx, bson.M{"istaken": true}, bson.M{"$set": bson.M{"istaken": false}})
	if err != nil {
		log.Errorf("Couldn't update taken bonuses: %s", err)
	}
}

func UpdateDailyBonus(ctx context.Context, b *domain.DailyBonus) error {
	upsert := true
	_, err := BonusCol().ReplaceOne(ctx, bson.M{"_id": b.UserID}, b, &options.ReplaceOptions{Upsert: &upsert})
	if err != nil {
		return err
	}
	return nil
}

func GetDailyBonus(ctx context.Context, userID string) (*domain.DailyBonus, error) {
	var b domain.DailyBonus
	err := BonusCol().FindOne(ctx, bson.M{"_id": userID}).Decode(&b)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func BonusCol() *mongo.Collection {
	// todo, rename to bonus
	return Client.Database(DbName).Collection("collectionName")
}
