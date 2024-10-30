package domain

import (
	"github.com/glossd/pokergloss/gomq/mqtable"
	"log"
)

type AssignmentType string

const (
	WinWithHand  AssignmentType = "winWithHand"
	LoseWithHand AssignmentType = "loseWithHand"

	WinWithPairOf  AssignmentType = "winWithPairOf"
	LoseWithPairOf AssignmentType = "loseWithPairOf"

	WinWithFace  AssignmentType = "winWithFace"
	LoseWithFace AssignmentType = "loseWithFace"

	BustPlayers   AssignmentType = "bustPlayers"
	DefeatPlayers AssignmentType = "defeatPlayers"
	ScareAway     AssignmentType = "scareAway"
	SharePot      AssignmentType = "sharePot"

	Win            AssignmentType = "win"
	WinLive        AssignmentType = "winLive"
	WinSitNGo      AssignmentType = "winSitNGo"
	WinMultiSitNGo AssignmentType = "winMultiSitNGo"
)

type AssignmentUnit struct {
	Type AssignmentType
	// Optional
	Hand
	Face   rune
	Number int
	TableRound
}

func New(t AssignmentType) *AssignmentUnit {
	return &AssignmentUnit{
		Type: t,
	}
}

func NewWithHand(t AssignmentType, hand Hand) *AssignmentUnit {
	if hand == "" {
		log.Panicf("Assignment with empty hand")
	}
	return &AssignmentUnit{
		Type: t,
		Hand: hand,
	}
}

func NewWithFace(t AssignmentType, pf rune) *AssignmentUnit {
	return &AssignmentUnit{
		Type: t,
		Face: pf,
	}
}

func NewWithNumber(t AssignmentType, n int) *AssignmentUnit {
	return &AssignmentUnit{
		Type:   t,
		Number: n,
	}
}

func NewWithNumberAndRound(t AssignmentType, n int, r TableRound) *AssignmentUnit {
	return &AssignmentUnit{
		Type:       t,
		Number:     n,
		TableRound: r,
	}
}

func (a *AssignmentUnit) getPrize() int64 {
	switch a.Type {
	case WinWithHand, LoseWithHand:
		switch a.Hand {
		case StraightFlush:
			return 50000
		case FourOfAKind:
			return 10000
		case FullHouse:
			return 2000
		case Flush:
			return 1500
		case Straight:
			return 1000
		case ThreeOfAKind:
			return 1000
		case TwoPair:
			return 750
		case Pair:
			return 500
		case HighCard:
			return 300
		default:
			return 0
		}
	case WinWithPairOf:
		return 1000
	case LoseWithPairOf:
		return 3000

	case WinWithFace:
		return 1000
	case LoseWithFace:
		return 1000

	case BustPlayers:
		switch a.Number {
		case 1:
			return 1000
		case 2:
			return 2000
		case 3:
			return 6000
		}
	case DefeatPlayers:
		switch a.Number {
		case 2:
			return 500
		case 3:
			return 1000
		case 4:
			return 2000
		}
	case Win:
		return 500
	case WinLive:
		return 500
	case WinSitNGo:
		return 750
	case WinMultiSitNGo:
		return 1500
	case ScareAway:
		switch a.Number {
		case 1:
			return 500
		case 2:
			return 1000
		case 3:
			return 2000
		case 4:
			return 3000
		}
	case SharePot:
		switch a.Number {
		case 1:
			return 1000
		case 2:
			return 2000
		}
	}
	return 0
}

func (a *AssignmentUnit) getMaxCount() int64 {
	switch a.Type {
	case WinWithHand, LoseWithHand:
		switch a.Hand {
		case StraightFlush, FourOfAKind:
			return 1
		default:
			return 3
		}
	case WinWithPairOf, LoseWithPairOf:
		return 1
	case LoseWithFace, WinWithFace:
		return 5
	case BustPlayers:
		switch a.Number {
		case 1:
			return 4
		case 2:
			return 2
		case 3:
			return 1
		}
	case DefeatPlayers:
		switch a.Number {
		case 2:
			return 4
		case 3:
			return 2
		case 4:
			return 1
		}
	case Win:
		return 10
	case WinLive:
		return 10
	case WinSitNGo:
		return 3
	case WinMultiSitNGo:
		return 3
	case ScareAway:
		switch a.Number {
		case 4:
			return 2
		case 3:
			return 3
		default:
			return 4
		}
	case SharePot:
		switch a.Number {
		case 2:
			return 2
		}
	}
	return 3
}

func (a *AssignmentUnit) matchGameEnd(p *mqtable.Player, ge *mqtable.GameEnd) bool {
	switch a.Type {
	case WinWithHand:
		return p.IsWinner && Hand(p.Hand) == a.Hand
	case LoseWithHand:
		return !p.IsWinner && Hand(p.Hand) == a.Hand
	case WinWithPairOf:
		return p.IsWinner && bothHasFace(p.Cards, a.Face)
	case LoseWithPairOf:
		return !p.IsWinner && bothHasFace(p.Cards, a.Face)

	case WinWithFace:
		return p.IsWinner && anyHasFace(p.Cards, a.Face)
	case LoseWithFace:
		return !p.IsWinner && anyHasFace(p.Cards, a.Face)

	case BustPlayers:
		if len(ge.Winners) == 1 {
			if p.IsWinner {
				var allInLostPlayers []*mqtable.Player
				for _, p := range ge.LostEndPlayers() {
					if p.LastAction == "allIn" {
						allInLostPlayers = append(allInLostPlayers, p)
					}
				}
				if len(allInLostPlayers) >= a.Number {
					return true
				}
			}
		}
	case DefeatPlayers:
		if len(ge.Winners) == 1 {
			if p.IsWinner {
				if len(ge.LostEndPlayers()) >= a.Number {
					return true
				}
			}
		}
	case Win:
		return p.IsWinner
	case WinLive:
		return p.IsWinner && ge.TableType == mqtable.GameEnd_LIVE
	case ScareAway:
		if !ge.GotToShowDown {
			if len(ge.Winners) == 1 {
				if p.IsWinner {
					lost := ge.LostChipsPlayers()
					if len(lost) > 0 {
						if a.TableRound != "" {
							if toTableRound(ge.TableRound) != a.TableRound {
								return false
							}
						}
						return len(lost) >= a.Number
					}
				}
			}
		}
	case SharePot:
		if p.IsWinner {
			if len(ge.Winners) > 1 {
				var sameHandCount int
				for _, winner := range ge.Winners {
					if winner.UserId != p.UserId {
						if winner.Hand == p.Hand {
							sameHandCount++
						}
					}
				}
				if sameHandCount >= a.Number {
					return true
				}
			}
		}
	}
	return false
}

func (a *AssignmentUnit) matchTournamentEnd(p *mqtable.Player, ge *mqtable.TournamentEnd) bool {
	switch a.Type {
	case WinSitNGo:
		for _, winner := range ge.TournamentWinners {
			if winner.UserId == p.UserId {
				return true
			}
		}
	case WinMultiSitNGo:
		for _, winner := range ge.TournamentWinners {
			if winner.UserId == p.UserId {
				return true
			}
		}
	}
	return false
}

func (a *AssignmentUnit) IsToLose() bool {
	return a.Type == LoseWithHand || a.Type == LoseWithPairOf || a.Type == LoseWithFace
}

func bothHasFace(cards []string, face rune) bool {
	if len(cards) < 2 {
		return false
	}
	for _, card := range cards {
		if len(card) != 2 {
			return false
		}
	}
	for _, card := range cards {
		if []rune(card)[0] != face {
			return false
		}
	}
	return true
}

func anyHasFace(cards []string, face rune) bool {
	if len(cards) < 2 {
		return false
	}
	for _, card := range cards {
		if len(card) != 2 {
			return false
		}
	}
	for _, card := range cards {
		if []rune(card)[0] == face {
			return true
		}
	}
	return false
}

func toTableRound(r mqtable.GameEnd_TableRound) TableRound {
	switch r {
	case mqtable.GameEnd_PRE_FLOP:
		return PreFlop
	case mqtable.GameEnd_FLOP:
		return Flop
	case mqtable.GameEnd_TURN:
		return Turn
	case mqtable.GameEnd_RIVER:
		return River
	}
	return ""
}
