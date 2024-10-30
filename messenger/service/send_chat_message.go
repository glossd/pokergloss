package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/domain"
	"github.com/glossd/pokergloss/messenger/web/ws/wssend"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrUseChatMessage = E("chat already exists, use chat message")

func SendMessageToChat(ctx context.Context, iden authid.Identity, chatID string, text string) (*domain.Message, error) {
	chatOID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, err
	}
	chat, err := db.FindChat(ctx, chatOID)
	if err != nil {
		return nil, err
	}
	if !chat.ContainsParticipant(iden.UserId) {
		return nil, err
	}
	msg, err := domain.NewMessage(chat.ID, iden, text)
	if err != nil {
		return nil, err
	}
	err = db.InsertMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	for _, participant := range chat.ParticipantsExcept(iden.UserId) {
		err = wssend.SendNewMessageTo(participant, msg)
		if err != nil {
			if !conf.IsE2E() {
				err := mqws.Publish(&mqws.Message{
					EntityType: mqws.Message_USER,
					EntityId:   participant,
					Events:     []*mqws.Event{{Type: "newsMessengerNewMessage", Payload: gomq.M{"from": iden, "text": text, "chatId": chatID}.JSON()}},
				})
				if err != nil {
					log.Errorf("Failed to publish mqws message: %s", err)
				}
			}
		}
	}

	return msg, nil
}
