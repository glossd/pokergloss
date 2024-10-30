package wssend

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/messenger/domain"
	"github.com/glossd/pokergloss/messenger/web/model"
	"github.com/glossd/pokergloss/messenger/web/ws/wsstore"
)

type M map[string]interface{}

func SendNewMessageTo(userID string, msg *domain.Message) error {
	client, err := wsstore.GetClient(userID)
	if err != nil {
		return err
	}

	client.SendChan <- event(wsstore.NewMessage, M{"message": model.ToMessage(msg)})
	return nil
}

func SendStatusTo(userID, msgID, chatID string, status domain.MessageStatus) error {
	client, err := wsstore.GetClient(userID)
	if err != nil {
		return err
	}

	client.SendChan <- event(wsstore.MessageStatus, M{"message": model.ToMessageStatus(msgID, chatID, status)})
	return nil
}

func SendNewChatTo(userID string, chat *model.Chat) error {
	client, err := wsstore.GetClient(userID)
	if err != nil {
		return err
	}

	client.SendChan <- event(wsstore.NewChat, M{"chat": chat})
	return nil
}

func SendTypingTo(userID string, chatID string, whoTyping authid.Identity) error {
	client, err := wsstore.GetClient(userID)
	if err != nil {
		return err
	}

	client.SendChan <- event(wsstore.Typing, M{"user": whoTyping, "chatId": chatID})
	return nil
}

func event(eventType wsstore.MessengerEventType, payload M) *wsstore.MessengerEvent {
	return &wsstore.MessengerEvent{Type: eventType, Payload: payload}
}
