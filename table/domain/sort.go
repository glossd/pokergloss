package domain

import "sort"

// Not a real linked list.
// Represents sorted players positions starting from first decided player position.
type linkedList []int

func (l linkedList) indexByPosition(position int) int {
	for i, v := range l {
		if v == position {
			return i
		}
	}
	return -1
}

func (l linkedList) isBefore(p1 int, p2 int) bool {
	i1 := l.indexByPosition(p1)
	i2 := l.indexByPosition(p2)
	return i1 < i2
}


func (l linkedList) isAfter(p1 int, p2 int) bool {
	i1 := l.indexByPosition(p1)
	i2 := l.indexByPosition(p2)
	return i1 > i2
}

// If not found returns -1
func (l linkedList) next(currentPos int) int {
	for i, pos := range l {
		if pos == currentPos {
			if i == len(l)-1 {
				return -1
			}
			return l[i+1]
		}
	}
	return -1
}

type linkedListFilter func (p *Player) bool

func (t *Table) sortPositionsFrom(firstPosition int, filter linkedListFilter) linkedList {
	result := make([]int, 0, t.Size)
	for _, seat := range t.Seats[firstPosition:] {
		if seat.IsTaken() {
			p := seat.GetPlayer()
			if filter(p) {
				result = append(result, p.Position)
			}
		}
	}
	for _, seat := range t.Seats[0:firstPosition] {
		if seat.IsTaken() {
			p := seat.GetPlayer()
			if filter(p) {
				result = append(result, p.Position)
			}
		}
	}
	return result
}

// First with the smallest TotalRoundBet
func SortPlayersByTotalRoundBet(players []*Player) {
	if len(players) > 1 {
		sort.Slice(players, func(i, j int) bool {
			return players[i].TotalRoundBet < players[j].TotalRoundBet
		})
	}
}

// First with the biggest TotalRoundBet
func SortPlayersByTotalRoundBetDesc(players []*Player) {
	if len(players) > 1 {
		sort.Slice(players, func(i, j int) bool {
			return players[i].TotalRoundBet > players[j].TotalRoundBet
		})
	}
}

// First with the biggest rank
func SortPlayersByHandRank(players []*Player) {
	if len(players) > 1 {
		sort.Slice(players, func(i, j int) bool {
			return players[i].HandRank < players[j].HandRank
		})
	}
}

// First the poorest
func SortPlayersByStartGameStack(players []*Player) {
	if len(players) > 1 {
		sort.Slice(players, func(i, j int) bool {
			return players[i].StartGameStack < players[j].StartGameStack
		})
	}
}

// First the poorest
func SortPlayersByStack(players []*Player) {
	if len(players) > 1 {
		sort.Slice(players, func(i, j int) bool {
			return players[i].Stack < players[j].Stack
		})
	}
}
