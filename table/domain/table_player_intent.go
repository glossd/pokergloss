package domain

import (
	log "github.com/sirupsen/logrus"
"github.com/glossd/pokergloss/auth/authid"
)

var ErrIntentNotSomeoneBetBefore = E("nobody bet before, you can't use this intent")
var ErrIntentNotNobodyBetBefore = E("someone bet before, you can't use this intent")
var ErrNoIntentsAvailable = E("no intents available")
var ErrRaiseIntentMaxRoundBet = E("raise intent can't be less than max round bet")
var ErrBetIntentMaxRoundBet = E("bet intent must be more than max round bet")
var ErrRaiseIntentNotEnoughChips = E("raise intent can't be more that your stack")
var ErrBetIntentNotEnoughChips = E("bet intent can't be more that your stack")

func (t *Table) SetIntent(position int, iden authid.Identity, intent Intent) error {
	p, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return err
	}

	maxBet := t.MaxRoundBet()

	if !t.IsInPlay() || t.IsDeciding(p) || !p.IsDecidable() || len(t.DecidablePlayers()) < 2 {
		return ErrNoIntentsAvailable
	}

	if t.isPlayerMadeRoundAction(position) {
		return ErrNoIntentsAvailable
	}

	if maxBet == 0 && !isNobodyBetBeforeIntent(intent.Type) {
		return ErrIntentNotNobodyBetBefore
	}

	if maxBet > p.TotalRoundBet && !isBetBeforeIntent(intent.Type) {
		return ErrIntentNotSomeoneBetBefore
	}

	if intent.IsAggressive() {
		maxAllowedBet := t.CalcMaxAllowedBet(p)
		var betChips int64
		if intent.Type == AllInIntentType {
			betChips = p.Stack
		} else {
			betChips = intent.Chips
		}
		if betChips > maxAllowedBet {
			return E("max bet is %d", maxAllowedBet)
		}
	}

	if intent.Type == CallIntentType || intent.Type == CallAnyIntentType {
		if p.Stack <= maxBet {
			log.Warnf("%s %s with stack less than max round bet, should've been all-in", p.Identity, intent.Type)
			return ErrNotEnoughChips
		}
	}

	if intent.Chips == p.Stack {
		log.Warnf("%s intent with whole stack, should've been all-in", intent.Type)
		intent = AllInIntent
	}

	if intent.Type == RaiseIntentType {
		if intent.Chips < maxBet {
			return ErrRaiseIntentMaxRoundBet
		}
		if intent.Chips > p.Stack {
			return ErrRaiseIntentNotEnoughChips
		}
	}

	if intent.Type == BetIntentType {
		if intent.Chips < t.BigBlind {
			return ErrBetIntentMaxRoundBet
		}
		if intent.Chips > p.Stack {
			return ErrBetIntentNotEnoughChips
		}
	}

	p.setIntent(intent)

	return nil
}

func (t *Table) isPlayerMadeRoundAction(position int) bool {
	pos, err := t.lastAggressorOrFirstRoundPosition()
	if err != nil {
		return false
	}

	linkedList := t.sortPositionsFrom(pos, func(*Player) bool { return true })
	return linkedList.isBefore(position, t.DecidingPosition)
}

func (t *Table) RemoveIntent(position int, iden authid.Identity) error {
	p, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return err
	}

	p.removeIntent()
	return nil
}
