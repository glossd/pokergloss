package e2e

import (
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/ws/web/mq"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

func TestSendUserMessage(t *testing.T) {
	cleanUp(t)

	ws, closeWs := wsDial(t, mqws.Message_USER, defaultToken)
	defer closeWs()

	mq.SendMessageToWs(&mqws.Message{
		EntityType: mqws.Message_USER,
		EntityId:   defaultIdentity.UserId,
		Events:     []*mqws.Event{{Type: "update", Payload: `{"balance":123}`}},
	})
	_, msg, err := ws.ReadMessage()
	assert.Nil(t, err)
	assert.EqualValues(t, "update", gjson.GetBytes(msg, "0.type").Str)
}
