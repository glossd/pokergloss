package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/messenger/web/clients/profile"
	log "github.com/sirupsen/logrus"
)

func InnerSendMessage(ctx context.Context, fromUserID, toUserID, text string) error {
	ucl, err := findOrCreateUserChatList(ctx, fromUserID)
	if err != nil {
		return err
	}
	user, err := profile.FetchUser(ctx, fromUserID)
	if err != nil {
		return err
	}
	fromIden := authid.Identity{UserId: fromUserID, Username: user.Username, Picture: user.Picture}
	if chat, ok := ucl.U2UChats[toUserID]; ok {
		_, err := SendMessageToChat(ctx, fromIden, chat.ID.Hex(), text)
		if err != nil {
			log.Errorf("Failed to inner send message, send message to chat: %s", err)
			return gomq.WrapInAckableError(err)
		}
	} else {
		chat, err := CreateU2UChat(ctx, fromIden, toUserID)
		if err != nil {
			log.Errorf("Failed to inner send message, create chat: %s", err)
			return gomq.WrapInAckableError(err)
		}
		_, err = SendMessageToChat(ctx, fromIden, chat.ID, text)
		if err != nil {
			log.Errorf("Failed to inner send message, after creating chat to chat: %s", err)
			return gomq.WrapInAckableError(err)
		}
	}
	return nil
}
