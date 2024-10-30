package model

import (
	"github.com/glossd/pokergloss/messenger/domain"
)

var EmptyString = ""

type Message struct {
	ID        string                `json:"id"`
	ChatID    *string               `json:"chatId,omitempty"`
	UserID    *string               `json:"userId,omitempty"`
	Text      *string               `json:"text,omitempty"`
	Status    *domain.MessageStatus `json:"status,omitempty" enums:"sent,read"`
	UpdatedAt *int64                `json:"updatedAt,omitempty"`
}

func ToMessage(msg *domain.Message) *Message {
	if msg == nil {
		return nil
	}
	chatID := msg.ChatID.Hex()
	updatedAt := msg.UpdatedAt * 1000
	return &Message{
		ID:        msg.ID.Hex(),
		ChatID:    &chatID,
		UserID:    &msg.UserID,
		Text:      &msg.Text,
		Status:    &msg.Status,
		UpdatedAt: &updatedAt,
	}
}

func ToMessageStatus(msgID, chatID string, status domain.MessageStatus) *Message {
	return &Message{
		ID:     msgID,
		ChatID: &chatID,
		Status: &status,
	}
}
