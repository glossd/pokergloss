package domain

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func (t *Table) dealFlop() {
	first, second, third := Algo.RandomAvailableThreeCards(t)
	t.CommunityCards.setFlop(first, second, third)
}

func (t *Table) dealTurn() {
	t.CommunityCards.setTurn(Algo.RandomAvailableOneCard(t))
}

func (t *Table) dealRiver() {
	t.CommunityCards.setRiver(Algo.RandomAvailableOneCard(t))
}

func (t *Table) dealCommunityCards() {
	switch t.CommunityCards.RoundType() {
	case PreFlopRound:
		t.dealFlop()
	case FlopRound:
		t.dealTurn()
	case TurnRound:
		t.dealRiver()
	}
}

// Recursion
func (t *Table) dealAllCommunityCards() {
	if t.CommunityCards.IsFull() {
		return
	}
	t.dealCommunityCards()
	t.dealAllCommunityCards()
}

// Rules for cash game: https://poker.stackexchange.com/a/2724/7972
func (t *Table) isRoundEnd() bool {
	players := t.GamingPlayers()
	if len(players) == 0 {
		log.Errorf("isRoundEnd invoked with no players on table, tableID=%s", t.ID.Hex())
		return false
	}
	if len(players) == 1 {
		return true
	}
	// all round bets should be matched to end the round
	maxRoundBet := t.MaxRoundBet()
	for _, player := range players {
		// checks if all players made a betting action
		if player.IsDecidable() && player.LastRoundAction == "" {
			return false
		}

		if player.IsDecidable() && maxRoundBet != player.TotalRoundBet {
			return false
		}
	}
	return true
}

func (t *Table) IsNewRound() bool {
	return t.wasNewRound
}

func (t *Table) IsPreFlop() bool {
	return t.CommunityCards.Flop == nil
}

func (t *Table) RoundType() RoundType {
	return t.CommunityCards.RoundType()
}

func (t *Table) IsFlop() bool {
	return t.CommunityCards.Flop != nil && t.CommunityCards.Turn == nil
}

func (t *Table) IsTurn() bool {
	cc := t.CommunityCards
	return cc.Flop != nil && cc.Turn != nil && cc.River == nil
}

func (t *Table) IsRiver() bool {
	cc := t.CommunityCards
	return cc.Flop != nil && cc.Turn != nil && cc.River != nil
}

func (t *Table) newRound() {
	t.wasNewRound = true
	t.LastAggressorPosition = -1
	t.moveBetsToPot()
	t.dealCommunityCards()
	for _, player := range t.AllPlayers() {
		player.newRoundReset()
	}

	var nextP, err = t.nextDecidablePlayer(t.DealerPosition())
	if err != nil {
		log.Errorf("Finishing game in newRound : %s", err)
		t.startShowDown()
		return
	}
	t.setActionDecidingPlayerPlusTime(nextP, time.Second)
}

func (t *Table) moveBetsToPot() {
	allInPlayers := t.PlayersFilter(func(p *Player) bool {
		return p.IsPlaying() && p.LastRoundAction == AllIn
	})
	notAllInGamingPlayers := t.DecidablePlayers()

	foldedPlayers := t.RoundFoldedPlayers()
	playersFoldedChips := make([]int64, 0, len(foldedPlayers))
	for _, p := range foldedPlayers {
		playersFoldedChips = append(playersFoldedChips, p.TotalRoundBet)
	}

	if len(allInPlayers) == 0 {
		t.Pots.increaseLastPot(t.roundBetChips())
	} else {
		var newPots []Pot
		SortPlayersByTotalRoundBet(allInPlayers)
		for i, allInPlayer := range allInPlayers {
			if i > 0 && allInPlayer.TotalRoundBet == allInPlayers[i-1].TotalRoundBet {
				// multiple players made all-in with same amount of chips
				// warning seems to be false
				newPots[len(newPots)-1].UserIDs = append(newPots[len(newPots)-1].UserIDs, allInPlayer.UserId)
				continue
			}

			var sumOfFoldedChips int64
			var newPlayersFoldedChips []int64
			for _, foldedChips := range playersFoldedChips {
				if foldedChips > allInPlayer.TotalRoundBet {
					sumOfFoldedChips += allInPlayer.TotalRoundBet
					newPlayersFoldedChips = append(newPlayersFoldedChips, foldedChips - allInPlayer.TotalRoundBet)
				} else {
					sumOfFoldedChips += foldedChips
				}
			}
			playersFoldedChips = newPlayersFoldedChips


			var prevAllInBet int64
			if i > 0 {
				prevAllInBet = allInPlayers[i-1].TotalRoundBet
			}

			betIncrease := allInPlayer.TotalRoundBet - prevAllInBet
			playersCount := int64(len(notAllInGamingPlayers) + len(allInPlayers) - i)
			pot := betIncrease*playersCount + sumOfFoldedChips
			newPots = append(newPots, Pot{Chips: pot, UserIDs: []string{allInPlayer.UserId}})
		}

		var chipsInNewPots int64
		for _, pot := range newPots {
			finishLastPot(&t.Pots, pot.Chips, pot.UserIDs)
			chipsInNewPots += pot.Chips
		}

		t.Pots.increaseLastPot(t.roundBetChips() - chipsInNewPots)
	}

	t.RoundPot = t.TotalPot
}

func (t *Table) roundBetChips() int64 {
	return t.TotalPot - t.RoundPot
}

func (t *Table) isLastRound() bool {
	return t.CommunityCards.RoundType() == RiverRound
}
