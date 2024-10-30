package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/domain"
	"github.com/glossd/pokergloss/messenger/web/clients/profile"
	"github.com/glossd/pokergloss/messenger/web/model"
	"github.com/glossd/pokergloss/messenger/web/ws/wssend"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateU2UChat(ctx context.Context, iden authid.Identity, withUserId string) (*model.Chat, error) {
	withUser, err := profile.FetchUser(ctx, withUserId)
	if err != nil {
		return nil, err
	}
	ucl, err := findOrCreateUserChatList(ctx, iden.UserId)
	if err != nil {
		return nil, err
	}
	if _, ok := ucl.U2UChats[withUserId]; ok {
		return nil, ErrUseChatMessage
	}

	chat := domain.NewChat(iden.UserId, withUserId)
	err = db.InsertChat(ctx, chat)
	if err != nil {
		return nil, err
	}
	ucl.SetChatWith(chat.ID, withUserId)
	err = db.SetUserListChatItem(ctx, ucl.UserID, ucl.U2UChats[withUserId])
	if err != nil {
		return nil, err
	}

	ulcOfToUser, err := db.FindUserChatList(ctx, withUserId)
	if err == mongo.ErrNoDocuments {
		ulcOfToUser = domain.NewUserChatList(withUserId)
		ulcOfToUser.SetChatWith(chat.ID, ucl.UserID)
		_ = db.InsertUserChatList(ctx, ulcOfToUser)
	} else if err != nil {
		log.Errorf("Failed to find user chat list of %s : %s", withUserId, err)
	} else {
		_ = db.SetUserListChatItem(ctx, withUserId, domain.NewChatForList(chat.ID, ucl.UserID))
	}

	_ = wssend.SendNewChatTo(withUserId, &model.Chat{
		ID:          chat.ID.Hex(),
		Name:        iden.Username,
		Picture:     iden.Picture,
		LastMessage: nil,
	})

	return &model.Chat{
		ID:          chat.ID.Hex(),
		Name:        withUser.Username,
		Picture:     withUser.Picture,
		LastMessage: nil,
	}, nil
}
