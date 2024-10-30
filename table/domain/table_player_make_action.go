package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	ErrRaiseMoreChips               = E("your raise requires to put more chips")
	ErrBetNotRaise                  = E("you can bet, but not raise")
	ErrRaiseButAllOtherPlayersAllIn = E("you can't raise, all other players made all-in")
	ErrAllInButAllOtherPlayersAllIn = E("you can't all-in, all other players made all-in")
)

func (t *Table) WasActionTimeout() bool {
	return t.wasTimeout
}

// Deprecated
func (t *Table) MakeActionDeprecated(position int, action ActionType, chips int64) error {
	player, err := t.GetPlayer(position)
	if err != nil {
		return err
	}
	return t.MakeAction(position, player.Identity, Action{Type: action, Chips: chips})
}

// Call, AllIn: auto fills chips,
func (t *Table) MakeAction(position int, iden authid.Identity, action Action) error {
	p, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return err
	}

	return t.makeActionWithValidation(p, action)
}

func (t *Table) makeActionWithValidation(p *Player, action Action) error {
	position := p.Position
	if !t.IsDeciding(p) {
		log.Warnf("User tried to make action decision out of turn, position=%d identity=%s", p.Position, p.Identity)
		return ErrNotYourTurn
	}

	t.madeActionPlayerPositions = nil
	t.wasTimeout = false
	maxRoundBet := t.MaxRoundBet()

	if action.Chips == 0 && action.IsHandBet() {
		switch action.Type {
		case Bet:
			action.Chips = t.BigBlind
		case Raise:
			action.Chips = maxRoundBet*2 - p.TotalRoundBet
		default:
			log.Errorf("User tried to put zero chips on %s action, identity=%s", action, p.Identity)
			return E("you must put chips on %s action", action)
		}
	}

	if action.Chips != 0 && action.IsChipFree() {
		log.Warnf("User tried to put number of chips on %s action, identity=%s", action, p.Identity)
		action.Chips = 0
	}

	if action.IsAggressive() {
		max := t.CalcMaxAllowedBet(p)
		if action.Type == AllIn {
			action.Chips = p.Stack
		}
		if action.Chips > max {
			return E("max bet is %d", max)
		}
	}

	// Call and AllIn mutate chips variable
	switch action.Type {
	case Check:
		if maxRoundBet > 0 && maxRoundBet != p.TotalRoundBet {
			log.Errorf("User tried to check but someone has bet before him, tableID=%s, maxRoundBet=%d, player=%v",
				t.ID.Hex(), maxRoundBet, p)
			return E("you can't check, because someone has a bigger bet")
		}
	case Call:
		if maxRoundBet == 0 {
			log.Errorf("User tried to call, but nobody called before him, tableID=%s, iden=%s", t.ID, p.Identity)
			return E("you can't call, nobody bet before you")
		}

		if maxRoundBet <= p.TotalRoundBet {
			log.Errorf("User tried to call, but his bet is more or equal to max round bet, tableID=%s, iden=%s", t.ID, p.Identity)
			return E("you can't call, your bet is more or equal to max round bet")
		}

	case Bet:
		if action.Chips == p.Stack {
			log.Warnf("User tried to bet with all stack he got, should've used allIn, tableID=%s, iden=%s", t.ID.Hex(), p.Identity)
			action.Type = AllIn
			break
		}

		if maxRoundBet != 0 {
			log.Errorf("User tried to bet, but max round bet is more than zero, tableID=%s, iden=%s", t.ID, p.Identity)
			return E("you can't bet, use raise")
		}
		if action.Chips < t.BigBlind {
			log.Errorf("User tried to bet with less chips than big blind, tableID=%s, iden=%s", t.ID, p.Identity)
			return E("your bet can't be less than big blind")
		}
	case Raise:
		if action.Chips == 0 {
			action.Chips = maxRoundBet*2 - p.TotalRoundBet
		}

		if action.Chips == p.Stack {
			log.Warnf("User tried to raise with all stack he got, should've used allIn, tableID=%s, iden=%s", t.ID.Hex(), p.Identity)
			action.Type = AllIn
			break
		}

		if maxRoundBet == 0 {
			log.Errorf("User tried to raise but no one bet, tableID=%s", t.ID.Hex())
			return E("you can't raise since no one bet, tableID=%s, iden=%s", t.ID, p.Identity)
		}

		if action.Chips <= (maxRoundBet - p.TotalRoundBet) {
			log.Errorf("User tried to raise to less then max round bet, tableID=%s, iden=%s", t.ID.Hex(), p.Identity)
			return ErrRaiseMoreChips
		}

		if p.Stack > maxRoundBet && t.isAllOtherPlayersAllIn(position) {
			log.Warnf("User tried raise when all others made all-in, tableID=%s, iden=%s", t.ID.Hex(), p.Identity)
			return ErrRaiseButAllOtherPlayersAllIn
		}
	case AllIn:
		if p.Stack > maxRoundBet && t.isAllOtherPlayersAllIn(position) {
			log.Warnf("User tried to all-in when all others made all-in, tableID=%s, iden=%s", t.ID.Hex(), p.Identity)
			return ErrAllInButAllOtherPlayersAllIn
		}
	}

	return t.doMakeAction(p, action)
}

// Only for project-internal use.
func (t *Table) MakeActionOnTimeout(position int) error {
	t.wasTimeout = true
	t.madeActionPlayerPositions = nil
	p, err := t.GetPlayer(position)
	if err != nil {
		log.Errorf("Couldn't make action on timeout: %s", err)
		return err
	}
	if !t.IsDeciding(p) {
		return ErrNotYourTurn
	}

	switch t.Status {
	case ShowdownTable:
		return t.makeShowdownActionOnTimeout(p)
	default:
		return t.makeBettingActionOnTimeout(p)
	}
}

func (t *Table) makeShowdownActionOnTimeout(p *Player) error {
	return t.MakeShowDownAction(Muck, p.Position, p.Identity)
}

func (t *Table) makeBettingActionOnTimeout(p *Player) error {
	action := Fold // ??? check or fold
	p.setToSittingOut()
	return t.doMakeAction(p, Action{action, 0})
}

// Makes player action without any validation.
func (t *Table) doMakeAction(p *Player, action Action) error {
	t.wasNewRound = false

	oldMaxRoundBet := t.MaxRoundBet()
	if action.Type == Call {
		// maxNum for case where bb is forced to all-in with less than bb chips
		action.Chips = maxNum(t.BigBlind-p.TotalRoundBet, oldMaxRoundBet-p.TotalRoundBet)
	}

	if action.Type != AllIn && action.Chips >= p.Stack {
		log.Warnf("Action has more chips than stack. Making auto allIn, tableID=%s, action=%+v, iden=%s", t.ID.Hex(), action, p.Identity)
		action = AllInAction
	}
	if action.Type == AllIn {
		action.Chips = p.Stack
	}

	if action.HasChips() {
		err := t.bet(p, action.Chips)
		if err != nil {
			return err
		}
		newMaxRoundBet := p.TotalRoundBet
		if newMaxRoundBet > oldMaxRoundBet {
			t.LastAggressorPosition = p.Position
			for _, player := range t.PlayersWithIntents() {
				player.updateIntent(newMaxRoundBet, oldMaxRoundBet)
			}
		}
	}

	p.makeChipsFreeAction(action)
	t.madeActionPlayerPositions = append(t.madeActionPlayerPositions, p)

	if t.shouldEndGame() {
		t.startShowDown()
		return nil
	}

	if t.isRoundEnd() {
		t.newRound()
		return nil
	}

	nextP, err := t.nextPlayer(p.Position, t.nextPlayerToDecideFilter())
	if err != nil {
		log.Fatalf("doMakeAction: didn't find next decidable player, table=%+v", t)
		return err
	}

	switch t.Type {
	case SitngoType, MultiType:
		if nextP.IsSittingOut() {
			return t.doMakeAction(nextP, FoldAction)
		}
		if nextP.IsLeaving {
			nextP.setToSittingOut()
			return t.doMakeAction(nextP, FoldAction)
		}
	case CashType:
		if nextP.IsLeaving {
			t.nullifyPlayer(nextP)
			return t.doMakeAction(nextP, FoldAction)
		}
	}

	if nextP.HasIntent() {
		nextIntentAction := nextP.GetIntentActionAndDelete()
		if nextIntentAction.Chips > nextP.Stack {
			log.Errorf("Next intent action has more chips than stack. Making auto allIn, tableID=%s, nextIntentAction=%+v, iden=%s", t.ID.Hex(), nextIntentAction, nextP.Identity)
		}
		return t.doMakeAction(nextP, nextIntentAction)
	}

	t.setActionDecidingPlayer(nextP)
	return nil
}

func (t *Table) setActionDecidingPlayer(p *Player) {
	t.setActionDecidingPlayerPlusTime(p, 0)
}

func (t *Table) setActionDecidingPlayerPlusTime(p *Player, plusTime time.Duration) {
	t.setToDecidingNoTimeout(p)
	duration := t.DecisionTimeout + plusTime + timeutil.Multiply(t.DecisionTimeout, p.additionalDecisionTimePercent())
	t.DecisionTimeoutAt = timeutil.NowAdd(duration)
}

func (t *Table) isAllOtherPlayersAllIn(position int) bool {
	players := t.PlayersFilter(func(p *Player) bool {
		return p.IsAllIn() && p.Position != position
	})
	if len(players) == (len(t.PlayingPlayersByGameType()) - 1) {
		return true
	}
	return false
}

func (t *Table) nextPlayerToDecideFilter() PlayerFilter {
	if t.IsTournament() {
		return func(p *Player) bool {
			if p.LastGameAction == "" && p.IsSittingOut() {
				return true
			}
			return p.IsDecidable()
		}
	}
	return DecidablePlayerFilter
}

func maxNum(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
