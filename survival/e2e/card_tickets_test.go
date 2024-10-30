package e2e

import (
	"context"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/domain"
	"github.com/glossd/pokergloss/survival/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFirstIncTicket(t *testing.T) {
	t.Cleanup(cleanUpDB)
	ctx := context.Background()
	assert.Nil(t, db.CardIncTicket(ctx, defaultIdentity.UserId))
	assert.EqualValues(t, 1, findCard(t).Tickets)
}

func TestTicketsDec(t *testing.T) {
	t.Cleanup(cleanUpDB)
	ctx := context.Background()
	assert.Nil(t, db.CardIncTicket(ctx, defaultIdentity.UserId))
	assert.Nil(t, db.CardDecTicket(ctx, defaultIdentity.UserId))
	assert.EqualValues(t, 0, findCard(t).Tickets)
	assert.NotNil(t, db.CardDecTicket(ctx, defaultIdentity.UserId))
	assert.EqualValues(t, 0, findCard(t).Tickets)
}

func TestTicketsDecService(t *testing.T) {
	t.Cleanup(cleanUpDB)
	ctx := context.Background()
	assert.Nil(t, db.CardIncTicket(ctx, defaultIdentity.UserId))

	dec, err := service.DecCardTickets(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.True(t, dec)

	dec, err = service.DecCardTickets(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.False(t, dec)
}

func findCard(t *testing.T) *domain.Card {
	card, err := db.FindCard(context.Background(), defaultIdentity.UserId)
	assert.Nil(t, err)
	return card
}
