package domain

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Assignment struct {
	AssignmentUnit
	Count int64
}

func NewAssignment(unit AssignmentUnit) *Assignment {
	return &Assignment{
		AssignmentUnit: unit,
		Count:          rand.Int63n(unit.getMaxCount())+1,
	}
}

func (a *Assignment) GetPrize() int64 {
	return a.getPrize()* a.Count
}

func (a *Assignment) GetFullName() string {
	pre, v, suf := a.GetName()
	return fmt.Sprintf("%s %v %s", pre, v, suf)
}

func (a *Assignment) GetName() (prefix string, variable interface{}, suffix string) {
	switch a.Type {
	case WinWithHand:
		return "Win with hand", a.Hand, ""
	case LoseWithHand:
		return "Lose with hand", a.Hand, ""
	case WinWithPairOf:
		return "Win with Pair of", string(a.Face), "in hands"
	case LoseWithPairOf:
		return "Lose with Pair of", string(a.Face), "in hands"
	case WinWithFace:
		return "Win with", string(a.Face), "in hands"
	case LoseWithFace:
		return "Lose with", string(a.Face), "in hands"
	case BustPlayers:
		return "Bust", a.Number, "players"
	case DefeatPlayers:
		return "Defeat", a.Number, "players"
	case Win:
		return "Win", "", ""
	case WinLive:
		return "Win", "", "Live"
	case WinSitNGo:
		return "Win SitNGo", "", ""
	case WinMultiSitNGo:
		return "Win Multi SitNGo", "", ""
	case ScareAway:
		middle := strconv.Itoa(a.Number)
		if a.Number == 0 {
			middle = ""
		}
		if a.TableRound != "" {
			middle = fmt.Sprintf("(%s) %s", a.TableRound.GetName(), middle)
		}
		return "Scare away", middle, "players"
	case SharePot:
		return "Share pot with", a.Number, "players"
	default:
		return "", "", ""
	}
}


