package domain

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func (t *Table) StartNextGame() error {
	if (t.IsMultiType() || t.IsSurvival) && t.IsWaiting() {
		return t.startFirstGame()
	}

	if !t.IsGameEnd() {
		log.Errorf("Couldn't start next game, table status is not %s but %s, tableID=%s", GameEndTable, t.Status, t.ID.Hex())
		return E("can't start next game, table status is not %s", GameEndTable)
	}

	dealerPos := t.DealerPosition()
	t.reset()

	t.handlePlayersWithZeroStack()
	t.PlayersCount = len(t.AllPlayers())

	if t.IsSurvival {
		if t.IsSurvivalUserLeft() {
			t.setToWaiting()
			return nil
		}
		// don't remove the user
		if t.PlayersCount == 1 {
			t.setToWaiting()
			return nil
		}
	}

	var newDealerPosition int
	switch t.Type {
	case CashType:
		t.nullifyStandingPlayers()
		t.PlayersCount = len(t.AllPlayers())
		if !t.IsEnoughPlayersForGame() {
			t.setToWaiting() // todo race condition
			return nil
		}
		p, err := t.nextPlayer(dealerPos, ActivePlayerFilter)
		if err != nil {
			log.Errorf("Couldn't start next game, didn't find next player, tableID=%s", t.ID.Hex())
			return err
		}
		newDealerPosition = p.Position
	case SitngoType, MultiType:
		if t.PlayersCount == 0 {
			t.setToWaiting()
			return nil
		}

		if t.PlayersCount == 1 {
			p, _ := t.IsOnlyOneOnTable()
			if t.IsSitngoType() {
				t.leaveSitngo(p)
			}
			if t.IsMultiType() && t.MultiAttrs.IsLast {
				p.tournamentInfo.IsLast = true
				t.nullifyPlayer(p)
			}
			t.setToWaiting()
			return nil
		}

		if t.IsSitngoType() {
			if t.AreAllSittingOut() {
				players := t.AllPlayers()
				SortPlayersByStack(players)
				for _, p := range players {
					t.leaveSitngo(p)
				}
				t.setToWaiting()
				return nil
			}
		}
		if t.IsMultiType() {
			if t.IsLast && t.AreAllSittingOut() {
				for _, p := range t.AllPlayers() {
					t.nullifyPlayer(p)
				}
				t.setToWaiting()
				return nil
			}
		}

		t.checkTimeAndIncreaseBlinds()

		p, err := t.nextPlayer(dealerPos, AnyPlayerFilter)
		if err != nil {
			log.Errorf("Table.StartNextGame failed, no nextPlayer after dealer: %s", err)
			return err
		}

		newDealerPosition = p.Position
	}

	return t.startNextGame(newDealerPosition)
}

func (t *Table) startFirstGame() error {
	var seatsToChoose []*Seat
	switch t.Type {
	case CashType:
		seatsToChoose = t.SeatsFilterByPlayerStatus(PlayerReady)
	case SitngoType, MultiType:
		seatsToChoose = t.AllTakenSeats()
	}

	dealer := Algo.ChooseDealer(seatsToChoose)
	return t.startNextGame(dealer.Position)
}

func (t *Table) startNextGame(newDealerPosition int) error {
	t.wasNewRound = true // should preFlop be a new round?
	t.Status = PlayingTable
	t.setReadyPlayersToPlaying()
	err := t.setBlinds(newDealerPosition)
	if err != nil {
		return err
	}
	t.dealHoleCardsToPlayers()

	if len(t.PlayingPlayersByGameType()) == 2 {
		// case when sb forced to bet allIn on blind.
		// case when bb forced to bet allIn with LessOrEqual to sb chips.
		// with two players on table the game shall end right away
		bbP, sbP, _ := t.BlindsPlayers()
		if bbP != nil && bbP.LastGameAction == AllIn && bbP.TotalRoundBet <= t.SmallBlind {
			t.isAutoGameEnd = true
			t.startShowDown()
			return nil
		}
		if sbP != nil && sbP.LastGameAction == AllIn {
			t.isAutoGameEnd = true
			t.startShowDown()
			return nil
		}
	}

	playerToDecide, err := t.nextPlayer(t.BigBlindPosition(), t.nextPlayerToDecideFilter())
	if err != nil {
		log.Errorf("Table#StartNextGame: player to decide wasn't found: %s", err)
		// not sure if I should end the game here
		t.startShowDown()
		return nil
	}
	t.setActionDecidingPlayerPlusTime(playerToDecide, time.Second)

	return t.postStart()
}

func (t *Table) setBlinds(newDealerPosition int) error {
	dealer := t.GetSeatUnsafe(newDealerPosition)
	dealer.setBlind(Dealer)

	var sb *Seat
	if len(t.PlayersFilter(t.GameTypePlayerFilter())) == 2 {
		sb = dealer
		sb.setBlind(DealerSmallBlind)
	} else {
		nextP, err := t.nextPlayer(dealer.Position, t.GameTypePlayerFilter())
		if err != nil {
			log.Errorf("setBlinds failed, no next playing player for sb, dPos=%d table=%+v", dealer.Position, t)
			return err
		}
		sb = t.GetSeatUnsafe(nextP.Position)
		sb.setBlind(SmallBlind)
	}

	t.betBlind(sb.GetPlayer())

	bb, err := t.nextPlayer(sb.Position, t.GameTypePlayerFilter())
	if err != nil {
		log.Errorf("setBlinds failed, no next playing player for bb, sbPos=%d table=%+v", sb.Position, t)
		return err
	}
	t.GetSeatUnsafe(bb.Position).setBlind(BigBlind)
	t.betBlind(bb)
	return nil
}

func (t *Table) dealHoleCardsToPlayers() {
	players := t.PlayingPlayersByGameType()

	for _, player := range players {
		firstCard, secondCard := Algo.RandomAvailableTwoCards(t)
		player.setHoleCards(firstCard, secondCard)
	}
}

func (t *Table) setReadyPlayersToPlaying() {
	seats := t.SeatsFilterByPlayerStatus(PlayerReady)
	for _, seat := range seats {
		seat.Player.Status = PlayerPlaying
	}
}

func (t *Table) reset() {
	t.TotalPot = 0
	t.RoundPot = 0
	t.Pots = initPots()
	t.Winners = nil
	t.CommunityCards.reset()
	t.DecidingPosition = -1
	t.LastAggressorPosition = -1

	for _, seat := range t.Seats {
		seat.reset()
	}
}

func (t *Table) handlePlayersWithZeroStack() {
	zeroStackPlayers := t.PlayersFilter(func(p *Player) bool {
		return p.Stack == 0
	})

	SortPlayersByStartGameStack(zeroStackPlayers)

	for _, p := range zeroStackPlayers {
		switch t.Type {
		case CashType:
			p.Status = PlayerReservedSeat
			t.brokePlayers = append(t.brokePlayers, p)
		case SitngoType:
			t.leaveSitngo(p)
		case MultiType:
			t.nullifyPlayer(p)
		}
	}
}

func (t *Table) PlayingPlayersByGameType() []*Player {
	return t.PlayersFilter(t.GameTypePlayerFilter())
}

func (t *Table) AreAllSittingOut() bool {
	all := t.AllPlayers()
	if len(all) == 0 {
		return false
	}
	var sittingOut int
	for _, player := range all {
		if player.IsSittingOut() {
			sittingOut++
		} else {
			return false
		}
	}
	return len(all) == sittingOut
}

func (t *Table) GameTypePlayerFilter() PlayerFilter {
	switch t.Type {
	case SitngoType, MultiType:
		return AnyPlayerFilter
	case CashType:
		return PlayingPlayerFilter
	}
	return PlayingPlayerFilter
}

func (t *Table) postStart() error {
	switch t.Type {
	case SitngoType, MultiType:
		decidingP, err := t.DecidingPlayer()
		if err == nil {
			if decidingP.IsSittingOut() {
				return t.doMakeAction(decidingP, FoldAction)
			}
		}
	}

	return nil
}
