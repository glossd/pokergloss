package domain

import (
	log "github.com/sirupsen/logrus"
)

const preFlopUserActionsLen = 10
const postFlopUserActionsLen = 20
var lastUserAction Action
var preFlopUserActions = make([]Action, 0, preFlopUserActionsLen)
var postFlopUserActions = make([]Action, 0, postFlopUserActionsLen)

// pipelines only for PostFlop
func (b *Bot) userTrackingPipeline(a Action) Action {
	if a == Fold && lastUserAction.Chips > 0 {
		if isUserPostflopAggressive() && b.confidence >= 0.1 {
			b.confidence += userPostflopAggro()
			return b.aggressionPipeline(b.loosenessPipeline())
		}
	}

	return a
}

func isUserPreflopAggressive() bool {
	aggro := userPreflopAggro()
	if aggro >= 0.7 {
		return len(preFlopUserActions) >= 3
	}
	if aggro > 0.5 {
		return len(preFlopUserActions) >= 5
	}
	return len(preFlopUserActions) > 10 && aggro > 0.25
}

func isUserPostflopAggressive() bool {
	aggro := userPostflopAggro()
	if aggro >= 0.7 {
		return len(postFlopUserActions) >= 3
	}
	if aggro > 0.5 {
		return len(postFlopUserActions) >= 5
	}
	return len(postFlopUserActions) > 10 && aggro > 0.25
}

func getLastPostflopUserActions(num int) []Action {
	if num == 0 {
		return nil
	}
	if len(postFlopUserActions) < num {
		return postFlopUserActions
	}
	return postFlopUserActions[len(postFlopUserActions)-num:]
}

func anyAggressiveAction(actions []Action) bool {
	for _, a := range actions {
		if a.Chips > 0 {
			return true
		}
	}
	return false
}

func userPreflopAggro() float64 {
	var aggro = make([]Action, 0, preFlopUserActionsLen)
	for _, action := range preFlopUserActions {
		if action.Chips > 0 {
			aggro = append(aggro, action)
		}
	}
	return float64(len(aggro))/float64(len(preFlopUserActions))
}

func userPostflopAggro() float64 {
	var aggro = make([]Action, 0, postFlopUserActionsLen/2)
	for _, action := range postFlopUserActions {
		if action.Chips > 0 {
			aggro = append(aggro, action)
		}
	}
	return float64(len(aggro))/float64(len(postFlopUserActions))
}

func UpdateUserTracker(oldT *Table, p *Player) {
	var action Action
	oldUserPlayer := oldT.Seats[p.Position].Player
	if p.LastGameActionType.IsAggressive() {
		action = Action{Type: p.LastGameActionType, Chips: oldUserPlayer.Stack - p.Stack}
	} else if p.LastGameActionType == AllInType {
		raisedChips := (oldUserPlayer.Stack + oldUserPlayer.TotalRoundBet) - oldT.MaxRoundBet
		if raisedChips <= 0 {
			action = AllIn
		} else {
			action = Action{Type: AllInType, Chips: raisedChips}
		}
	} else {
		action = Action{Type: p.LastGameActionType}
	}

	if oldT.IsPreFlop() {
		if len(preFlopUserActions) >= preFlopUserActionsLen {
			preFlopUserActions = append(preFlopUserActions[1:], action)
		} else {
			preFlopUserActions = append(preFlopUserActions, action)
		}
	} else {
		if len(postFlopUserActions) >= postFlopUserActionsLen {
			postFlopUserActions = append(postFlopUserActions[1:], action)
		} else {
			postFlopUserActions = append(postFlopUserActions, action)
		}
	}

	lastUserAction = action
	if oldT.IsPreFlop() {
		log.Debugf("User PreFlop Aggro: %.2f", userPreflopAggro())
	} else {
		log.Debugf("User PostFlop Aggro: %.2f", userPostflopAggro())
	}
}
