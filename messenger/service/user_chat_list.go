package service

import (
	"context"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

func findOrCreateUserChatList(ctx context.Context, userID string) (*domain.UserChatList, error) {
	ucl, err := db.FindUserChatList(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ucl = domain.NewUserChatList(userID)
			err := db.InsertUserChatList(ctx, ucl)
			if err != nil {
				return nil, err
			}
			return ucl, nil
		} else {
			return nil, err
		}
	}
	return ucl, nil
}

func findOrBuildUserChatList(ctx context.Context, userID string) (*domain.UserChatList, error) {
	ucl, err := db.FindUserChatList(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.NewUserChatList(userID), nil
		} else {
			return nil, err
		}
	}
	return ucl, nil
}
