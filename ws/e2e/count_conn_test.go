package e2e

import (
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/ws/storage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserConnCount(t *testing.T) {
	cleanUp(t)

	assert.EqualValues(t, 0, storage.GetUserConnections())

	ws, closeWs := wsDial(t, mqws.Message_USER, defaultToken)
	defer closeWs()

	time.Sleep(time.Millisecond) // Hub.Register is async operation
	assert.EqualValues(t, 1, storage.GetUserConnections())

	err := ws.Close()
	assert.Nil(t, err)

	time.Sleep(time.Millisecond)
	assert.EqualValues(t, 0, storage.GetUserConnections())
}
