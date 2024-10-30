package db

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/bank/domain"
	log "github.com/sirupsen/logrus"
	"github.com/tevino/abool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const listLimit = 20

var ErrPageLimit = fmt.Errorf("max page limit is %d", listLimit)

// true: firstView, false: secondView
var viewSwitch = abool.New()

const firstView = "firstRatings"
const secondView = "secondRatings"

func viewName() string {
	if viewSwitch.IsSet() {
		return firstView
	} else {
		return secondView
	}
}

func viewNameOpposite() string {
	if viewSwitch.IsSet() {
		return secondView
	} else {
		return firstView
	}
}

func BuildOppositeRankView() {
	ctx, cancelDrop := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancelDrop()
	err := rankView(viewNameOpposite).Drop(ctx)
	if err != nil {
		log.Errorf("Couldn't drop opposite rank view, %s", err)
		return
	}

	pipeline := bson.A{
		bson.M{"$match": bson.M{"nonratable": false}},
		bson.M{"$sort": bson.D{{"chips", -1}, {"updatedat", -1}}},
		bson.D{{"$group", bson.D{
			{"_id", nil},
			{"items", bson.D{{"$push", bson.M{"_id": "$_id", "username": "$username", "picture": "$picture", "chips": "$chips", "updatedat": "$updatedat"}}}},
		}}},
		bson.M{"$unwind": bson.M{"path": "$items", "includeArrayIndex": "rank"}},
		bson.M{"$group": bson.D{
			{"_id", "$items._id"},
			{"username", bson.M{"$first": "$items.username"}},
			{"picture", bson.M{"$first": "$items.picture"}},
			{"chips", bson.M{"$sum": "$items.chips"}},
			{"updatedat", bson.M{"$sum": "$items.updatedat"}},
			{"rank", bson.M{"$sum": bson.M{"$add": bson.A{"$rank", 1}}}}, // rank + 1
		}},
		bson.M{"$sort": bson.D{{"chips", -1}, {"updatedat", -1}}},
	}

	ctx, cancelCreate := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCreate()
	err = Client.Database(DbName).CreateView(ctx, viewNameOpposite(), BalanceCol().Name(), pipeline)
	if err != nil {
		log.Errorf("Couldn't create rank view, %s", err)
		return
	}

	log.Debugf("Built opposite rank view, switching views")
	viewSwitch.Toggle()
}

func FindRatings(ctx context.Context, skip, limit int64) ([]*domain.Rating, error) {
	if limit > 20 {
		return nil, ErrPageLimit
	}
	cur, err := RankView().Find(ctx, bson.M{}, &options.FindOptions{Skip: &skip, Limit: &limit})
	if err != nil {
		return nil, err
	}
	var results []*domain.Rating
	for cur.Next(ctx) {
		var r domain.Rating
		err := cur.Decode(&r)
		if err != nil {
			return nil, err
		}
		results = append(results, &r)
	}
	return results, nil
}

func FindRatingsNoCtx(skip, limit int64) ([]*domain.Rating, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindRatings(ctx, skip, limit)
}

func FindRating(ctx context.Context, userID string) (*domain.Rating, error) {
	var rank domain.Rating
	err := RankView().FindOne(ctx, bson.D{{"_id", userID}}).Decode(&rank)
	if err != nil {
		return nil, err
	}

	return &rank, nil
}

func FindRatingNoCtx(userID string) (*domain.Rating, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindRating(ctx, userID)
}

func CountRatings(ctx context.Context) (int64, error) {
	return RankView().CountDocuments(ctx, bson.D{})
}

func rankView(name func() string) *mongo.Collection {
	return Client.Database(DbName).Collection(name())
}

func RankView() *mongo.Collection {
	return rankView(viewName)
}
