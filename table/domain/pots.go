package domain

import log "github.com/sirupsen/logrus"

// Zero index pot is the main pot and others are side pots
type Pots []*Pot

type Pot struct {
	idx int
	Chips int64
	// all-in players who share this pot
	UserIDs []string
	// winners who get this pot
	WinnerPositions []int
}

func (pots Pots) getPlayerPotsIdx(userId string) int {
	for potIdx, pot := range pots {
		for _, uid := range pot.UserIDs {
			if uid == userId {
				return potIdx
			}
		}
	}
	return len(pots) - 1
}

func (pots Pots) isWinnerOfPot(position int) bool {
	for _, pot := range pots {
		for _, pos := range pot.WinnerPositions {
			if pos == position {
				return true
			}
		}
	}
	return false
}

// from included
// to included
func (pots Pots) sumOfPots(from, to int) int64 {
	if from > to {
		return 0
	}
	var sum int64
	for i := from; i <= to; i++ {
		if i == -1 {
			continue
		}
		sum += pots[i].Chips
	}
	return sum
}

func (pots Pots) slice(from int) []*Pot {
	if from == -1 {
		from = 0
	}
	if from >= len(pots) {
		return nil
	}
	return pots[from:]
}

func (pots Pots) setIndexes() {
	for i, pot := range pots {
		pot.idx = i
	}
}

func (t *Table) removeLastPotIfEmpty() {
	lastPot :=  t.Pots[len(t.Pots)-1]
	if lastPot.Chips == 0 {
		t.Pots = t.Pots[:len(t.Pots)-1]
	}
}

func (t *Table) buildWinners() (winners []Winner) {
	positionChips := make(map[int]int64)
	var orderedPositions []int
	for _, pot := range t.Pots {
		for _, position := range pot.WinnerPositions {
			if _, ok := positionChips[position]; !ok {
				orderedPositions = append(orderedPositions, position)
			}
			positionChips[position] += pot.Chips/int64(len(pot.WinnerPositions))
		}
	}
	for _, position := range orderedPositions {
		p, err := t.GetPlayer(position)
		if err != nil {
			log.Errorf("domain.Table.buildWinners player position %d not found ", position)
			continue
		}
		winners = append(winners, Winner{Position: position, Chips: positionChips[position], HandRank: p.HandRankString})
	}
	return
}


func (pots Pots) setWinnerPos(from, to, winnerPos int) {
	for i := from; i <= to; i++ {
		pot := pots[i]
		if !pot.containsWinnerPos(winnerPos) {
			// todo somehow same winnerPos can be many times here
			pot.WinnerPositions = append(pot.WinnerPositions, winnerPos)
		}
	}
}

func (pot *Pot) containsWinnerPos(winnerPos int) bool {
	if pot == nil {
		return false
	}
	for _, position := range pot.WinnerPositions {
		if position == winnerPos {
			return true
		}
	}
	return false
}

func (pots Pots) setAllWinnerPos(winnerPos int) {
	wp := []int{winnerPos}
	for _, pot := range pots {
		pot.WinnerPositions = wp
	}
}

func (pots Pots) increaseLastPot(chips int64) {
	pots[len(pots)-1].Chips += chips
}

func finishLastPot(pots *Pots, chips int64, userIDs []string) {
	potsVal := *pots
	potsVal[len(potsVal)-1].Chips += chips
	potsVal[len(potsVal)-1].UserIDs = userIDs
	*pots = append(potsVal, &Pot{})
}

func initPots() Pots {
	return Pots{{}}
}
