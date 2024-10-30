package actionhandler

import (
	"github.com/glossd/pokergloss/table/domain"
)

func GetActionHandler(table *domain.Table) ActionHandler {
	if table.IsGameEnd() {
		return GameEnd{}
	}

	if table.IsShowDown() {
		return ShowDown{}
	}

	if table.IsNewRound() {
		return RoundEnd{}
	}

	return Action{}
}
