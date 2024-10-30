package db

import (
	"context"
	"github.com/glossd/pokergloss/profile/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var True = true

func UpsertProfileNoCtx(p *domain.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return UpsertProfile(ctx, p)
}

func SaveProfile(ctx context.Context, p *domain.Profile) error {
	_, err := ProfileCol().InsertOne(ctx, p)
	if err != nil {
		log.Errorf("Couldn't save user to mongo: %s", err)
		return err
	}
	return nil
}

func UpsertProfile(ctx context.Context, p *domain.Profile) error {
	_, err := ProfileCol().ReplaceOne(ctx, bson.D{{Key: "_id", Value: p.Username}}, p, &options.ReplaceOptions{Upsert: &True})
	if err != nil {
		log.Errorf("Couldn't save user to mongo: %s", err)
		return err
	}
	return nil
}

func UpdateUserID(username, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err := ProfileCol().UpdateOne(ctx, bson.D{{"_id", username}}, bson.D{{"$set", bson.D{{"userid", userID}}}})
	return err
}

func UpdatePicture(username, pictureURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err := ProfileCol().UpdateOne(ctx, bson.D{{"_id", username}}, bson.D{{"$set", bson.D{{"picture", pictureURL}}}})
	return err
}

func DeleteProfile(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err := ProfileCol().DeleteOne(ctx, bson.D{{"_id", username}})
	return err
}

func FindProfile(ctx context.Context, username string) (*domain.Profile, error) {
	var doc domain.Profile
	err := ProfileCol().FindOne(ctx, bson.D{{"_id", username}}).Decode(&doc)
	if err != nil {
		log.Errorf("Failed to findOne usernameDoc: %s", err)
		return nil, err
	}
	return &doc, err
}

func FindAllProfilesByUserIDs(ctx context.Context, userIDs []string) ([]*domain.Profile, error) {
	cur, err := ProfileCol().Find(ctx, bson.D{{"userid", bson.M{"$in": userIDs}}})
	if err != nil {
		return nil, err
	}

	var profiles []*domain.Profile
	for cur.Next(ctx) {
		var doc domain.Profile
		err := cur.Decode(&doc)
		if err != nil {
			return profiles, err
		}
		profiles = append(profiles, &doc)
	}

	if err := cur.Err(); err != nil {
		log.Warnf("Find profiles by userIds failed: %s", err)
		return profiles, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(profiles) == 0 {
		return []*domain.Profile{}, nil
	}

	return profiles, nil
}

var SearchLimit int64 = 20

func SearchProfiles(ctx context.Context, username string) ([]*domain.Profile, error) {
	var profiles []*domain.Profile
	cur, err := ProfileCol().Find(ctx, bson.M{"_id": primitive.Regex{Pattern: username, Options: "i"}}, &options.FindOptions{Limit: &SearchLimit})
	if err != nil {
		log.Errorf("Search profiles failed: %s", err)
		return nil, err
	}

	for cur.Next(ctx) {
		var doc domain.Profile
		err := cur.Decode(&doc)
		if err != nil {
			return profiles, err
		}
		profiles = append(profiles, &doc)
	}

	if err := cur.Err(); err != nil {
		log.Errorf("Search profiles failed: %s", err)
		return profiles, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(profiles) == 0 {
		return []*domain.Profile{}, nil
	}

	return profiles, nil
}

func FindProfileNoCtx(username string) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return FindProfile(ctx, username)
}

func ExistsUsername(ctx context.Context, username string) (bool, error) {
	findOptions := options.Find()
	findOptions.SetLimit(1)
	cur, err := ProfileCol().Find(ctx, bson.D{{"_id", username}}, findOptions)
	if err != nil {
		return false, err
	}
	return cur.Next(ctx), nil
}

func ProfileCol() *mongo.Collection {
	return Client.Database(DbName).Collection("usernames")
}
