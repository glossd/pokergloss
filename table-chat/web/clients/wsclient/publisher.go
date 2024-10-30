package wsclient

import (
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table-chat/domain"
	log "github.com/sirupsen/logrus"
)

var TestMQ = make(chan *mqws.TableMessage, 64)

type TableChatEvent struct {
	Text string          `json:"text"`
	User authid.Identity `json:"user"`
}

// @ID use websocket instead
// @Summary table chat events
// @Success 200 {object} TableChatEvent
// @Router /use-websocket [get]
func PublishChatMessage(chatMsg *domain.Message) error {
	msg := &mqws.TableMessage{
		ToEntityIds: []string{chatMsg.TableID},
		Events:      []*mqws.Event{{Type: "chatMessage", Payload: gomq.M{"text": chatMsg.Text, "user": chatMsg.CreatedBy}.JSON()}},
	}
	if conf.IsProd() {
		err := mqws.PublishTableMsg(msg)
		if err != nil {
			log.Errorf("Failed to publish chat message: %s", err)
		}
		return err
	}

	if conf.IsE2E() {
		TestMQ <- msg
	}

	return nil
}

func PublishEmojiMessage(emojiMsg *domain.EmojiMessage) error {
	msg := &mqws.TableMessage{
		ToEntityIds: []string{emojiMsg.TableID},
		Events:      []*mqws.Event{{Type: "emojiMessage", Payload: gomq.M{"emoji": emojiMsg.Emoji, "user": emojiMsg.CreatedBy}.JSON()}},
	}
	if conf.IsProd() {
		err := mqws.PublishTableMsg(msg)
		if err != nil {
			log.Errorf("Failed to publish chat message: %s", err)
		}
		return err
	}

	if conf.IsE2E() {
		TestMQ <- msg
	}

	return nil
}
