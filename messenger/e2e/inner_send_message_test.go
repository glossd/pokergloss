package e2e

import (
	"context"
	"github.com/glossd/pokergloss/messenger/service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInnerSendMessage(t *testing.T) {
	t.Cleanup(cleanUp)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := service.InnerSendMessage(ctx, "1", "2", "Hi")
	assert.Nil(t, err)
}
