package playerbank

import (
	"fmt"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/web/client/bankclient"
)

func SendPlayerChipsToBank(p *domain.Player, t *domain.Table) {
	if t.IsSurvival {
		return
	}
	switch t.Type {
	case domain.CashType:
		if p.Stack != 0 {
			bankclient.Deposit(p.Stack, p.UserId, fmt.Sprintf("Left table %s", t.Name))
		}
	case domain.SitngoType:
		if p.GetTournamentInfo().Prize > 0 {
			bankclient.Deposit(p.GetTournamentInfo().Prize, p.UserId, fmt.Sprintf("Won Sit & Go %s", t.Name))
		}
	case domain.MultiType:
		if p.GetTournamentInfo().Prize > 0 {
			bankclient.Deposit(p.GetTournamentInfo().Prize, p.UserId, fmt.Sprintf("Won Multi Table %s", t.Name))
		}
	}
}
