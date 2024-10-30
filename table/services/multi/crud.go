package multi

import (
	"context"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/paging"
)

func FindLobbies(ctx context.Context, params paging.Params) ([]*domain.LobbyMulti, error) {
	return db.FindMultiLobbiesSortedByStartAt(ctx, params)
}
