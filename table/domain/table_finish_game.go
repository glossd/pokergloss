package domain

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
)

func (t *Table) finishGame() {
	t.previousDecidingPosition = t.DecidingPosition
	t.setToGameEnd()
	t.rake = t.buildRake()
	// uses the rake
	t.potsToPlayersStack()

	t.DecidingPosition = -1
	timeoutAcc := conf.Props.GameEndMinTimeout
	for range t.Pots {
		timeoutAcc += conf.Props.GameEndPotTimeout
	}
	if t.CommunityCards.isNewCardsAreDealtOnStartShowDown {
		newCardsLen := len(t.CommunityCards.GetNewCards())
		switch newCardsLen {
		case 1:
			timeoutAcc += conf.Props.GameEndCommunityCardTimeout
		case 2:
			timeoutAcc += 2 * conf.Props.GameEndCommunityCardTimeout
		case 5:
			timeoutAcc += 3 * conf.Props.GameEndCommunityCardTimeout
		}
	}
	t.DecisionTimeoutAt = timeutil.NowAdd(timeoutAcc)
}

func (t *Table) potsToPlayersStack() {
	rake := t.GetRake()
	for _, w := range t.Winners {
		chipsToStack := w.Chips - rake.Of(w.Position)
		t.GetPlayerUnsafe(w.Position).addWinningsToStack(chipsToStack)
	}
}

func (t *Table) isReadyForGame() bool {
	if t.Status == WaitingTable {
		if len(t.SeatsFilterByPlayerStatus(PlayerReady)) >= 2 {
			return true
		}
	}
	return false
}

func (t *Table) IsEnoughPlayersForGame() bool {
	return len(t.ActivePlayers()) >= 2
}

func (t *Table) IsEnoughPlayersForNextGame() bool {
	playersForNextGame := t.PlayersFilter(func(p *Player) bool {
		return p.IsAvailableForNextGame()
	})
	return len(playersForNextGame) >= 2
}

func (t *Table) shouldEndGame() bool {
	if len(t.GamingPlayers()) <= 1 {
		return true
	}

	if t.isRoundEnd() && t.isLastRound() {
		return true
	}

	decidablePlayers := t.DecidablePlayers()
	if len(decidablePlayers) < 1 {
		return true
	}
	if len(decidablePlayers) == 1 {
		maxRoundBet := t.MaxRoundBet()
		oneDecidableP := decidablePlayers[0]
		if oneDecidableP.TotalRoundBet >= maxRoundBet {
			return true
		}
	}

	return false
}

func (t *Table) isOneDecidableLeft() bool {
	return len(t.DecidablePlayers()) == 1
}

func (t *Table) isOneGamingLeft() bool {
	return len(t.GamingPlayers()) == 1
}

func (t *Table) isOneGamingLeftPlayer() (*Player, bool) {
	players := t.GamingPlayers()
	if len(players) == 1 {
		return players[0], true
	}
	return nil, false
}

func (t *Table) PlayersMoreThan(count int) bool {
	players := t.AllPlayers()
	if len(players) > count {
		return true
	}
	return false
}

func (t *Table) IsOnlyOneOnTableBool() bool {
	_, only := t.IsOnlyOneOnTable()
	return only
}

func (t *Table) IsOnlyOneOnTable() (*Player, bool) {
	players := t.AllPlayers()
	if len(players) == 1 {
		return players[0], true
	}
	return nil, false
}

func (t *Table) IsOnlyOneOnTableReady() bool {
	players := t.PlayersFilter(func(p *Player) bool {
		return p.Status == PlayerReady
	})

	if len(players) == 1 {
		return true
	}

	return false
}

func (t *Table) lastGamingPlayer() *Player {
	players := t.GamingPlayers()
	if len(players) != 1 {
		log.Errorf("Player is not last one, tableID=%s", t.ID)
		return nil
	}
	return players[0]
}

func (t *Table) IsNewGameStarted() bool {
	return t.IsNewRound() && t.IsPreFlop()
}

// https://en.wikipedia.org/wiki/Showdown_(poker)
func (t *Table) WasShowDown() bool {
	players := t.PlayersFilter(func(p *Player) bool {
		return p.Status == PlayerPlaying && p.LastGameAction != Fold
	})
	if len(players) > 1 {
		return true
	}
	return false
}

func (t *Table) GetBigBlind() *Player {
	for _, seat := range t.Seats {
		if seat.IsTaken() {
			player := seat.Player
			if player.Blind == BigBlind {
				return player
			}
		}
	}
	return nil
}

func (t *Table) Blinds() (int, int, int) {
	bb, sb, d := -1, -1, -1
	for _, seat := range t.Seats {
		switch seat.Blind {
		case BigBlind:
			bb = seat.Position
		case SmallBlind:
			sb = seat.Position
		case Dealer:
			d = seat.Position
		case DealerSmallBlind:
			d = seat.Position
			sb = seat.Position
		}
	}
	if bb == -1 || sb == -1 || d == -1 {
		log.Warnf("Blind is not found in table=%s, %t, %t, %t", t.ID.Hex(), bb == -1, sb == -1, d == -1)
	}
	return bb, sb, d
}

// Returns seats with BigBlind, SmallBlind and Dealer positions
func (t *Table) BlindsPlayers() (bb *Player, sb *Player, d *Player) {
	for _, seat := range t.Seats {
		if seat.IsTaken() {
			player := seat.Player
			switch player.Blind {
			case BigBlind:
				bb = seat.Player
			case SmallBlind:
				sb = seat.Player
			case Dealer:
				d = seat.Player
			case DealerSmallBlind:
				d = seat.Player
				sb = seat.Player
			}
		}
	}
	return bb, sb, d
}

func (t *Table) BlindsSlice() []*Player {
	bb, sb, d := t.BlindsPlayers()
	return []*Player{bb, sb, d}
}

func (t *Table) UniqueBlindsSlice() []*Player {
	bb, sb, d := t.BlindsPlayers()
	if sb == d {
		return []*Player{bb, d}
	}
	return []*Player{bb, sb, d}
}

func (t *Table) BigBlindPosition() int {
	bb, _, _ := t.Blinds()
	return bb
}

func (t *Table) SmallBlindPosition() int {
	_, sb, _ := t.Blinds()
	return sb
}

func (t *Table) SmallBlindPlayer() *Player {
	_, sb, _ := t.BlindsPlayers()
	return sb
}

func (t *Table) DealerPosition() int {
	_, _, d := t.Blinds()
	return d
}

func (t *Table) DealerPlayer() *Player {
	_, _, d := t.BlindsPlayers()
	return d
}

// At least one player made AllIn
func (t *Table) IsAllInGame() bool {
	for _, seat := range t.Seats {
		if seat.IsTaken() {
			if seat.GetPlayer().LastGameAction == AllIn {
				return true
			}
		}
	}
	return false
}
