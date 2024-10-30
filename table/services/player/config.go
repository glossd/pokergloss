package player

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/playerautoconfig"
)

func SetAutoMuck(params *PositionParams, autoMuck bool) error {
	err := applyAutoConfig(params, func(t *domain.Table) error {
		return t.SetAutoMuck(autoMuck, params.position, params.iden)
	})
	if err != nil {
		return err
	}

	_ = playerautoconfig.SetAutoMuck(params.ctx, params.iden, autoMuck)

	return nil
}

func SetAutoTopUp(params *PositionParams, autoTopUp bool) error {
	err := applyAutoConfig(params, func(t *domain.Table) error {
		return t.SetAutoTopUp(autoTopUp, params.position, params.iden)
	})
	if err != nil {
		return err
	}

	_ = playerautoconfig.SetAutoTopUp(params.ctx, params.iden, autoTopUp)

	return nil
}

func SetAutoReBuy(params *PositionParams, autoReBuy bool) error {
	err := applyAutoConfig(params, func(t *domain.Table) error {
		return t.SetAutoReBuy(autoReBuy, params.position, params.iden)
	})
	if err != nil {
		return err
	}

	_ = playerautoconfig.SetAutoReBuy(params.ctx, params.iden, autoReBuy)

	return nil
}

func applyAutoConfig(params *PositionParams, apply func(*domain.Table) error) error {
	table, err := db.FindTable(params.ctx, params.tableID)
	if err != nil {
		return err
	}

	err = apply(table)
	if err != nil {
		return err
	}
	p, err := table.GetPlayer(params.position)
	if err != nil {
		return err
	}

	err = db.SetTableContext(params.ctx, table.ID, db.PlayerAutoConfig(p.Position, p.GetAutoConfig()))
	if err != nil {
		return nil
	}

	return nil
}
