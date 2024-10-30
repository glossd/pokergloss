package multi

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/bankclient"
	log "github.com/sirupsen/logrus"
)

func Unregister(params services.IdenParams) error {
	lobby, err := db.FindLobbyMulti(params.GetCtx(), params.GetID())
	if err != nil {
		return err
	}

	_, err = lobby.Unregister(params.GetIden())
	if err != nil {
		return err
	}

	err = db.UpdateLobbyMulti(params.GetCtx(), lobby)
	if err != nil {
		return err
	}

	err = bankclient.Deposit(lobby.BuyIn, params.GetIden().UserId, "Unregistered from tournament")
	if err != nil {
		log.Errorf("IMPORTANT! failed to send back %d to user %s", lobby.BuyIn, params.GetIden())
	}

	broadcast.SendTableEvent(params.GetID().Hex(), &events.TableEvent{
		Type:    events.MultiUnregisterType,
		Payload: events.M{"player": params.GetIden(), "prizes": lobby.Prizes, "prizePool": lobby.PrizePool()},
	})

	return nil
}
