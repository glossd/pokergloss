package e2e

import (
	"context"
	"github.com/glossd/pokergloss/assignment/db"
	"github.com/glossd/pokergloss/assignment/service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateDaily(t *testing.T) {
	t.Cleanup(cleanUp)
	now := time.Date(2020, 8, 6, 0, 0, 0, 0, time.Now().Location())
	_, err := service.CreateDaily(now)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	d, err := db.FindDaily(ctx, "2020-08-06")
	assert.Nil(t, err)
	assert.EqualValues(t, 3, len(d.Assignments))
}
