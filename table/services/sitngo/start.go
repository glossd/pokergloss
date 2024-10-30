package sitngo

import (
	"context"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
)

func Start(ctx context.Context, startAt int64) error {
	lobbies, err := db.FindSitAndGoLobbiesForStartAt(ctx, startAt)
	if err != nil {
		return err
	}
	for _, lobby := range lobbies {
		err := lobby.StartAnyway()
		if err != nil {
			if err != domain.ErrSitngoNotEnoughPlayers {
				log.Errorf("Failed to start sitngo by startAt: %s", err)
			}
			continue
		}
		err = db.UpdateSitAndGoLobby(ctx, lobby)
		if err != nil {
			return err
		}
		handleTable(ctx, lobby)
	}
	return nil
}
