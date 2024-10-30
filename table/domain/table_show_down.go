package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
	"math"
)

type ShowDownActionType string

const (
	Show      ShowDownActionType = "show"
	Muck      ShowDownActionType = "muck"
	ShowLeft  ShowDownActionType = "showLeft"
	ShowRight ShowDownActionType = "showRight"
)

func (t *Table) startShowDown() {

	if !t.CommunityCards.IsFull() && t.IsAllInGame() && !t.isOneGamingLeft() {
		t.dealAllCommunityCards()
		t.CommunityCards.isNewCardsAreDealtOnStartShowDown = true
	}
	t.moveStackOverflow()
	t.moveBetsToPot()
	t.ComputeWinners()

	position, err := t.lastAggressorOrFirstRoundPosition()
	if err != nil {
		t.finishGame()
		return
	}

	t.setToShowDown()

	t.showDown(position, true)
}

// include posFrom to Show Down
func (t *Table) showDown(posFrom int, include bool) {
	firstPositionToShowDown, err := t.lastAggressorOrFirstRoundPosition()
	if err != nil {
		t.finishGame()
		return
	}

	if p, onlyOne := t.isOneGamingLeftPlayer(); onlyOne {
		if p.ShowDownAction != "" {
			t.finishGame()
			return
		}
		if p.AutoConfig.Muck {
			t.setShowDownAction(p, Muck)
			t.finishGame()
			return
		} else {
			t.setShowDownDecidingPlayer(p)
			return
		}
	}

	if len(t.Winners) == 0 {
		log.Errorf("Zero winners on showdown, tableID=%s", t.ID)
		t.finishGame()
		return
	}

	sortedPos := t.sortPositionsFrom(firstPositionToShowDown, func(p *Player) bool { return p.IsGaming() })
	startFromPos := posFrom
	if !include {
		startFromPos = sortedPos.next(startFromPos)
		p, err := t.GetPlayer(startFromPos)
		if err != nil || p.ShowDownAction != "" {
			t.finishGame()
			return
		}

	}
	var startFromPosIdx int
	for i, pos := range sortedPos {
		if pos == startFromPos {
			startFromPosIdx = i
			break
		}
	}

	lastBestHandRank := int32(math.MaxInt32)
	for _, pos := range sortedPos[startFromPosIdx:] {
		p, err := t.GetPlayer(pos)
		if err != nil {
			log.Errorf("Error to start show down, position with no player: %s", err)
			continue
		}
		if t.CommunityCards.isNewCardsAreDealtOnStartShowDown {
			// https://poker.stackexchange.com/a/2709/7972
			// situation where every player gone all-in before river
			t.setShowDownAction(p, Show)
			continue
		}

		// vip showdown privilege
		if p.IsVip() {
			if !t.isWinner(p.Position) {
				if p.AutoConfig.Muck {
					t.setShowDownAction(p, Muck)
					continue
				} else {
					t.setShowDownDecidingPlayer(p)
					return
				}
			}
		}

		if p.HandRank <= lastBestHandRank {
			lastBestHandRank = p.HandRank
			t.setShowDownAction(p, Show)
		} else {
			if t.isWinner(p.Position) {
				t.setShowDownAction(p, Show)
			} else {
				if p.AutoConfig.Muck {
					t.setShowDownAction(p, Muck)
				} else {
					t.setShowDownDecidingPlayer(p)
					return
				}
			}
		}
	}

	t.finishGame()
}

// Case when a player made a bet more then max all-in of players,
// or everybody made fold after aggression
func (t *Table) moveStackOverflow() {
	players := t.PlayingPlayersByGameType()
	if len(players) > 1 {
		SortPlayersByTotalRoundBetDesc(players)
		firstBetPlayer := players[0]
		secondBetPlayer := players[1]
		chipsToMove := firstBetPlayer.TotalRoundBet - secondBetPlayer.TotalRoundBet
		if chipsToMove > 0 {
			err := firstBetPlayer.moveChipsFromBetToStack(chipsToMove)
			if err != nil {
				log.Errorf("Table.moveStackOverflow failed: %s", err)
			}
			t.TotalPot -= chipsToMove
			t.stackOverflowPlayer = firstBetPlayer
		}
	}
}

func (t *Table) setShowDownAction(p *Player, action ShowDownActionType) {
	p.setShowDownAction(action)
	t.showedDownPlayers = append(t.showedDownPlayers, p)
}

func (t *Table) setShowDownDecidingPlayer(p *Player) {
	t.setToDecidingNoTimeout(p)
	t.DecisionTimeoutAt = timeutil.NowAdd(conf.Props.Table.ShowDownTimeout)
}

func (t *Table) MakeShowDownAction(action ShowDownActionType, position int, iden authid.Identity) error {
	t.showedDownPlayers = nil

	if !t.IsShowDown() {
		return E("table is not in showdown state")
	}

	p, err := t.GetPlayerIdentified(position, iden)
	if err != nil {
		return err
	}

	if t.DecidingPosition != position {
		return ErrNotYourTurn
	}

	p.setShowDownAction(action)
	t.showedDownPlayers = append(t.showedDownPlayers, p)

	t.showDown(position, false)
	return nil
}
