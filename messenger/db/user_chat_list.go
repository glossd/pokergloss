package db

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/messenger/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUserChatList(ctx context.Context, userID string) (*domain.UserChatList, error) {
	var ucl domain.UserChatList
	err := UserChatListCol().FindOne(ctx, filterID(userID)).Decode(&ucl)
	if err != nil {
		return nil, err
	}
	return &ucl, nil
}

func InsertUserChatList(ctx context.Context, list *domain.UserChatList) error {
	_, err := UserChatListCol().InsertOne(ctx, list)
	if err != nil {
		log.Errorf("Failed to insert user chat list of userID=%s: %s", list.UserID, err)
		return err
	}
	return nil
}

func SetUserListChatItem(ctx context.Context, userID string, chat *domain.U2UChatForList) error {
	var isSupportCreated bool
	if chat.OtherUserID == domain.SupportUserID {
		isSupportCreated = true
	}
	_, err := UserChatListCol().UpdateOne(ctx, filterID(userID), bson.M{"$set": bson.M{
		fmt.Sprintf("u2uchats.%s", chat.OtherUserID): chat,
		"issupportcreated":                           isSupportCreated,
	}})
	if err != nil {
		log.Errorf("Failed to set chats last message of userID=%s: %s", userID, err)
		return err
	}
	return nil
}

func UserChatListCol() *mongo.Collection {
	return Client.Database(DbName).Collection("userChatList")
}
