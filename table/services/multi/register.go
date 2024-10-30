package multi

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/bankclient"
	log "github.com/sirupsen/logrus"
)

func Register(params services.IdenParams) error {
	lobby, err := db.FindLobbyMulti(params.GetCtx(), params.GetID())
	if err != nil {
		return err
	}

	err = lobby.Register(params.GetIden())
	if err != nil {
		return err
	}

	err = bankclient.Withdraw(params.GetCtx(), lobby.BuyIn, params.GetIden().UserId, "Registered for tournament")
	if err != nil {
		return err
	}

	err = db.UpdateLobbyMulti(params.GetCtx(), lobby)
	if err != nil {
		err = bankclient.Deposit(lobby.BuyIn, params.GetIden().UserId, "Failed to register for tournament")
		if err != nil {
			log.Errorf("IMPORTANT! failed to send back %d to user %s", lobby.BuyIn, params.GetIden())
		}
		return err
	}

	broadcast.SendTableEvent(params.GetID().Hex(), &events.TableEvent{
		Type:    events.MultiRegisterType,
		Payload: events.M{"player": params.GetIden(), "prizes": lobby.Prizes, "prizePool": lobby.PrizePool()},
	})

	return nil
}
