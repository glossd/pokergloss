package service

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/table-chat/db"
	"github.com/glossd/pokergloss/table-chat/domain"
	"github.com/glossd/pokergloss/table-chat/web/clients/wsclient"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

var ErrEmpty = fmt.Errorf("text can't be empty")
var ErrInvalid = fmt.Errorf("data is invalid")

func PostMessage(ctx context.Context, tableID string, iden authid.Identity, inputText string) error {
	text := strings.TrimSpace(inputText)
	if text == "" {
		return ErrEmpty
	}
	if IsUserInBlackList(iden.UserId) {
		return ErrUserBlackListed
	}
	message := &domain.Message{
		ID:        primitive.NewObjectID(),
		TableID:   tableID,
		CreatedBy: iden,
		CreatedAt: time.Now().UnixNano() / 1e6,
		Text:      text,
	}

	db.InsertMessage(ctx, message)
	return wsclient.PublishChatMessage(message)
}

func PostEmoji(ctx context.Context, tableID string, iden authid.Identity, emoji string) error {
	if emoji == "" {
		return ErrEmpty
	}
	if !domain.IsValidEmoji(emoji) {
		return ErrInvalid
	}
	msg := &domain.EmojiMessage{
		ID:        primitive.NewObjectID(),
		TableID:   tableID,
		CreatedBy: iden,
		CreatedAt: time.Now().UnixNano() / 1e6,
		Emoji:     domain.Emoji(emoji),
	}
	_ = db.InsertEmojiMessage(ctx, msg)
	return wsclient.PublishEmojiMessage(msg)
}
