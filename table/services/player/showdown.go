package player

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
)

func MakeShowDownAction(params *PositionParams, action domain.ShowDownActionType) error {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		return err
	}

	err = table.MakeShowDownAction(action, params.position, params.iden)
	if err != nil {
		return err
	}

	err = actionhandler.Handle(params.ctx, table)
	if err != nil {
		return err
	}

	return nil
}

func GetConfig(params *PositionParams) (*domain.PlayerAutoConfig, error) {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		return nil, err
	}

	p, err := table.GetPlayerIdentified(params.position, params.iden)
	if err != nil {
		return nil, err
	}

	return p.GetAutoConfig(), nil
}
