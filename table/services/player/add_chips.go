package player

import (
	"fmt"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/bankclient"
	log "github.com/sirupsen/logrus"
)

func AddChips(params *ChipsParams) error {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		return err
	}

	return AddChipsOnTable(params, table)
}

func AddChipsOnTable(params *ChipsParams, table *domain.Table) error {
	err := table.AddChips(params.chips, params.position, params.iden)
	if err != nil {
		return err
	}
	p, err := table.GetPlayer(params.position)
	if err != nil {
		return err
	}

	err = bankclient.Withdraw(params.ctx, params.chips, params.iden.UserId, fmt.Sprintf("Added chips to table %s", table.Name))
	if err != nil {
		if err == bankclient.ErrNotEnoughChips {
			return err
		}
		log.Errorf("Couldn't add chips for %s, bank service error: %s", params.iden, err)
		return fmt.Errorf("bank service unavailable: %s", err)
	}
	if p.ChipsToAddOnReset > 0 {
		err = db.SetTable(params.tableID, db.PlayerAddChips(p))
		if err != nil {
			return err
		}
	} else {
		err = db.SetTableGameFlow(params.ctx, params.tableID, table.GameFlowVersion, db.PlayerAddChips(p))
		if err != nil {
			return err
		}
	}

	broadcast.SendTableEvent(params.tableID.Hex(), events.BuildAddChips(p))

	return nil
}
