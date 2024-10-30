package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var ErrYouCantChangeStatus = E("you can't change message status to 'read'")
var ErrMessageTextEmpty = E("message test can't be empty")

type MessageStatus string

const (
	Sent MessageStatus = "sent"
	Read MessageStatus = "read"
)

type Message struct {
	ID     primitive.ObjectID `bson:"_id"`
	ChatID primitive.ObjectID
	// the one who sent the message
	UserID    string
	Text      string
	Status    MessageStatus
	CreatedAt int64 // seconds
	UpdatedAt int64
}

func NewMessage(chatID primitive.ObjectID, from authid.Identity, text string) (*Message, error) {
	if text == "" {
		return nil, ErrMessageTextEmpty
	}
	now := time.Now().Unix()
	return &Message{
		ID:        primitive.NewObjectID(),
		ChatID:    chatID,
		UserID:    from.UserId,
		Text:      text,
		Status:    Sent,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
