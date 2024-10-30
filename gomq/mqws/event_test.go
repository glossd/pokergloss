package mqws

import (
	"github.com/glossd/pokergloss/gomq"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

func TestMarshalEvent(t *testing.T) {
	e1 := &Event{
		Type:    "balance",
		Payload: gomq.M{"balance": 123}.JSON(),
	}

	events := []*Event{e1}
	msg := Message{Events: events}
	result := msg.EventsToJson()

	assert.EqualValues(t, 1, gjson.GetBytes(result, "#").Int())
	assert.EqualValues(t, 123, gjson.GetBytes(result, "0.payload.balance").Int())
}
