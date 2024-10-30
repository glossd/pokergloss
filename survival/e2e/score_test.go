package e2e

import (
	"context"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoreFindDefault(t *testing.T) {
	t.Cleanup(cleanUpDB)

	orDefault := findScore(t)
	assert.EqualValues(t, orDefault.Level, 0)
}

func TestScoreUpsertMax(t *testing.T) {
	t.Cleanup(cleanUpDB)

	ctx := context.Background()

	assert.Nil(t, db.UpsertScore(ctx, defaultIdentity.UserId, 4))
	assert.EqualValues(t, 4, findScore(t).Level)

	assert.Nil(t, db.UpsertScore(ctx, defaultIdentity.UserId, 7))
	assert.EqualValues(t, 7, findScore(t).Level)

	assert.Nil(t, db.UpsertScore(ctx, defaultIdentity.UserId, 3))
	assert.EqualValues(t, 7, findScore(t).Level)
}

func findScore(t *testing.T) *domain.Score {
	orDefault, err := db.FindScoreOrDefault(context.Background(), defaultIdentity.UserId)
	assert.Nil(t, err)
	return orDefault
}
