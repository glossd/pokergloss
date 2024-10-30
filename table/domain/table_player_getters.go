package domain

import (
	log "github.com/sirupsen/logrus"
"github.com/glossd/pokergloss/auth/authid"
"math/rand"
)

var ErrNoPlayersAtTable = E("no players at table")

var AnyPlayerFilter = func(p *Player) bool { return true }
var PlayingPlayerFilter = func(p *Player) bool { return p.IsPlaying() }
var ActivePlayerFilter = func(p *Player) bool { return p.IsActive() }
var GamingPlayerFilter = func(p *Player) bool { return p.IsGaming() }
var DecidablePlayerFilter = func(p *Player) bool { return p.IsDecidable() }
var SittingOutPlayerFilter = func(p *Player) bool { return p.IsSittingOut() }
var BrokePlayerFilter = func(p *Player) bool { return p.IsBroke() }

func (t *Table) GetPlayerIdentified(position int, iden authid.Identity) (*Player, error) {
	seat, err := t.GetSeatIdentified(position, iden)
	if err != nil {
		return nil, err
	}
	if seat.IsFree() {
		return nil, E("no player is sitting at position %d", position)
	}
	return seat.Player, nil
}

func (t *Table) GetPlayer(position int) (*Player, error) {
	seat, err := t.GetSeat(position)
	if err != nil {
		return nil, err
	}
	if seat.IsFree() {
		return nil, E("seat is empty")
	}
	return seat.GetPlayer(), nil
}

func (t *Table) GetRandomPlayer() (*Player, error) {
	all := t.AllPlayers()
	if len(all) == 0 {
		return nil, ErrNoPlayersAtTable
	}

	return all[rand.Intn(len(all))], nil
}

func (t *Table) IsSeatFree(position int) bool {
	seat, err := t.GetSeat(position)
	if err != nil {
		return false
	}

	return seat.IsFree()
}

// Use it really wisely, maybe only in tests
func (t *Table) GetPlayerUnsafe(position int) *Player {
	return t.GetSeatUnsafe(position).Player
}

func (t *Table) DecidingPlayerUnsafe() *Player {
	p, err := t.DecidingPlayer()
	if err != nil {
		log.Errorf("Tried to get deciding player unsafe: %s", err)
	}
	return p
}

func (t *Table) DecidingPlayer() (*Player, error) {
	if t.DecidingPosition < 0 {
		return nil, E("deciding position is not set")
	}
	player, err := t.GetPlayer(t.DecidingPosition)
	if err != nil {
		return nil, E("deciding position doesn't have a player, decidingPosition=%d", t.DecidingPosition)
	}
	return player, nil
}

func (t *Table) PreviousDecidingPlayer() *Player {
	if t.previousDecidingPosition < 0 {
		log.Errorf("Tried to get previous deciding player, but previous deciding position is not set")
		return nil
	}
	player, err := t.GetPlayer(t.previousDecidingPosition)
	if err != nil {
		log.Errorf("Previous deciding position doesn't have a player, position=%d", t.previousDecidingPosition)
	}
	return player
}

func (t *Table) AllGamePlayers() []*Player {
	return t.PlayersFilter(func(p *Player) bool {
		return p.LastGameAction != ""
	})
}

func (t *Table) ActivePlayers() []*Player {
	return t.PlayersFilter(func(p *Player) bool { return p.IsActive() })
}

func (t *Table) AutoTopUpPlayers() []*Player {
	return t.PlayersFilter(func(p *Player) bool { return p.AutoConfig.TopUp && p.IsTopUpable() })
}

func (t *Table) AutoReBuyPlayers() []*Player {
	return t.PlayersFilter(func(p *Player) bool { return p.AutoConfig.ReBuy && p.IsReBuyable() })
}

func (t *Table) DecidablePlayers() []*Player {
	return t.PlayersFilter(func(p *Player) bool { return p.IsDecidable() })
}

func (t *Table) PlayersWithIntents() []*Player {
	return t.PlayersFilter(func(p *Player) bool { return p.IsDecidable() && p.HasIntent() })
}

func (t *Table) PlayersWithUpgradedIntents() []*Player {
	return t.PlayersFilter(func(p *Player) bool { return p.IsDecidable() && (p.isIntentRemoved || p.isIntentChanged) })
}

func (t *Table) GamingPlayers() []*Player {
	return t.PlayersFilter(func(p *Player) bool { return p.IsGaming() })
}

func (t *Table) RoundFoldedPlayers() []*Player {
	return t.PlayersFilter(func(p *Player) bool {
		return p.IsRoundFolded()
	})
}

func (t *Table) ReservedSeatPlayers() []*Player {
	return t.PlayersFilter(func(p *Player) bool {
		return p.Status == PlayerReservedSeat
	})
}

func (t *Table) AllPlayers() []*Player {
	return t.PlayersFilter(AnyPlayerFilter)
}

func (t *Table) nextDecidablePlayer(currentPos int) (*Player, error) {
	return t.nextPlayer(currentPos, func(p *Player) bool { return p.IsDecidable() })
}

func (t *Table) nextGamingPlayer(currentPos int) (*Player, error) {
	return t.nextPlayer(currentPos, func(p *Player) bool { return p.IsGaming() })
}

func (t *Table) nextActivePlayer(currentPos int) (*Player, error) {
	return t.nextPlayer(currentPos, func(p *Player) bool { return p.IsActive() })
}

func (t *Table) nextPlayingPlayerUnsafe(currentPosition int) *Player {
	p, _ := t.nextPlayer(currentPosition, func(player *Player) bool {
		return player.Status == PlayerPlaying
	})
	return p
}

func (t *Table) nextPlayingPlayer(currentPosition int) (*Player, error) {
	return t.nextPlayer(currentPosition, func(p *Player) bool { return p.IsPlaying() })
}

func (t *Table) nextPlayer(currentPosition int, filter PlayerFilter) (*Player, error) {
	err := t.validatePosition(currentPosition)
	if err != nil {
		return nil, err
	}
	allPlayers := t.AllPlayers()
	if len(allPlayers) == 0 {
		log.Warnf("Next player invoked with no players at the table, tableID=%s", t.ID.Hex())
		return nil, E("nextPlayer: no player")
	}
	if len(allPlayers) == 1 {
		log.Warnf("Next player invoked with one player at the table, tableID=%s", t.ID.Hex())
		return allPlayers[0], E("nextPlayer: only one player")
	}

	if currentPosition < len(t.Seats)-1 {
		seatsAfterCurrent := t.Seats[currentPosition+1:]
		for _, seat := range seatsAfterCurrent {
			if seat.IsTaken() {
				player := seat.GetPlayer()
				if filter(player) {
					return player, nil
				}
			}
		}
	}
	for _, seat := range t.Seats[:currentPosition] {
		if seat.IsTaken() {
			player := seat.GetPlayer()
			if filter(player) {
				return player, nil
			}
		}
	}

	return nil, E("no next player found")
}

type PlayerFilter func(*Player) bool

func (t *Table) PlayersFilter(filter PlayerFilter) []*Player {
	var players []*Player
	for _, seat := range t.Seats {
		if seat.IsTaken() {
			if filter(seat.GetPlayer()) {
				players = append(players, seat.GetPlayer())
			}
		}
	}
	return players
}

func (t *Table) PlayerSeatsFilter(filter PlayerFilter) []*Seat {
	var seats []*Seat
	for _, seat := range t.Seats {
		if seat.IsTaken() {
			if filter(seat.GetPlayer()) {
				seats = append(seats, seat)
			}
		}
	}
	return seats
}

func (t *Table) SeatsFilterByPlayerStatus(ps PlayerStatus) []*Seat {
	var seats []*Seat
	for _, seat := range t.Seats {
		if seat.Player != nil && seat.Player.Status == ps {
			seats = append(seats, seat)
		}
	}
	return seats
}

func (t *Table) AllTakenSeats() []*Seat {
	var seats []*Seat
	for _, seat := range t.Seats {
		if seat.IsTaken() {
			seats = append(seats, seat)
		}
	}
	return seats
}

func (t *Table) ContainsPlayer(iden authid.Identity) bool {
	return t.PositionOf(iden) > -1
}

func (t *Table) FindPlayerByIdentity(iden authid.Identity) *Player {
	playerPosition := t.PositionOf(iden)
	if playerPosition > -1 {
		return t.GetPlayerUnsafe(playerPosition)
	}
	return nil
}

// Returns position of user.
// If user is not present the method returns -1.
func (t *Table) PositionOf(iden authid.Identity) int {
	players := t.AllPlayers()
	for i, player := range players {
		if player.UserId == iden.UserId {
			return i
		}
	}
	return -1
}
