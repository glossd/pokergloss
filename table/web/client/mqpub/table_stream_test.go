package mqpub

import (
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToEvents(t *testing.T) {
	from := []*events.TableEvent{{Type: "first", Payload: events.M{"foo": "bar"}}}
	result := ToEvents(from)
	assert.Len(t, result, 1)
	assert.EqualValues(t, "first", result[0].Type)
}
