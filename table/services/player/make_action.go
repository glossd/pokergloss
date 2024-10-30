package player

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
)

func MakeBettingAction(params *ChipsParams, action domain.ActionType) error {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		return err
	}
	return DoBettingActionOnTable(table, params, action)
}

// Made to test race condition
func DoBettingActionOnTable(table *domain.Table, params *ChipsParams, action domain.ActionType) error {
	ctx := params.ctx
	iden := params.iden
	position := params.position
	chips := params.chips

	p, err := table.GetPlayer(position)
	if err == nil {
		allowedToActTime := p.UpdatedAt + timeutil.Duration(conf.Props.PlayerActionDuration)
		if timeutil.Now() < allowedToActTime {
			return services.ErrFormat("wait...")
		}
	}

	err = table.MakeAction(position, iden, domain.Action{Type: action, Chips: chips})
	if err != nil {
		return err
	}

	err = actionhandler.Handle(ctx, table)
	if err != nil {
		return err
	}

	return nil
}
