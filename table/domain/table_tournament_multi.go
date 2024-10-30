package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math/rand"
)

var ErrNoAvailableSeats = E("no available seats")
var ErrNotAvailableInMultiTable = E("no available in multi table game type")

type MultiAttrs struct {
	 PlayerMoves []PlayerMove
	 putPlayers  []*Player
	 movedPlayers map[primitive.ObjectID][]*Player
	 IsLast bool
}

func (m MultiAttrs) GetPutPlayers() []*Player {
	return m.putPlayers
}

func (m MultiAttrs) MovedPlayersPerTable() map[primitive.ObjectID][]*Player {
	return m.movedPlayers
}

type PlayerMove struct {
	FromPosition int
	ToTableID primitive.ObjectID
}

func NewTableMulti(l *LobbyMulti, params NewTableParams, seats []*Seat) (*Table, error) {
	t, err := NewTable(params)
	if err != nil {
		return nil, err
	}

	t.Type = MultiType

	t.setSeatsForTournament(seats)
	err = t.startFirstGame()
	if err != nil {
		return nil, err
	}

	t.TournamentAttributes = l.tournamentAttrs()

	return t, nil
}

func (t *Table) MultiPutPlayerAtFreePosition(p *Player) error {
	freePos := t.FirstAvailablePosition()
	if freePos < 0 {
		return ErrNoAvailableSeats
	}
	seat := t.GetSeatUnsafe(freePos)
	seat.multiSitPlayer(p)
	t.MultiAttrs.putPlayers = append(t.MultiAttrs.putPlayers, p)
	return nil
}

func (t *Table) MultiAddPlayerToMove(movingPlayer *Player, moveToTable *Table) {
	t.PlayerMoves = append(t.PlayerMoves,
		PlayerMove{
			FromPosition: movingPlayer.Position,
			ToTableID:    moveToTable.ID,
		})
	movingPlayer.isMultiMoving = true
}

func (t *Table) MultiDeletePlayerForMove(movingPlayerSeat *Seat, toTableID primitive.ObjectID) {
	if t.MultiAttrs.movedPlayers == nil {
		t.MultiAttrs.movedPlayers = make(map[primitive.ObjectID][]*Player)
	}
	t.MultiAttrs.movedPlayers[toTableID] = append(t.MultiAttrs.movedPlayers[toTableID], movingPlayerSeat.GetPlayer())
	movingPlayerSeat.RemovePlayer()
}

func (t *Table) MultiRandomAvailablePlayerToMove() *Player {
	allPlayers := t.AllPlayers()
	availablePlayers := make([]*Player, 0, len(allPlayers))
	for _, p := range allPlayers {
		if !p.isMultiMoving {
			availablePlayers = append(availablePlayers, p)
		}
	}
	if len(availablePlayers) == 0 {
		log.Panicf("No available player to move for multi, tableId=%s", t.ID.Hex())
	}
	return availablePlayers[rand.Intn(len(availablePlayers))]
}

func (t *Table) MultiSitPlayerAndTryToStartGame(p *Player) (bool, error) {
	err := t.MultiPutPlayerAtFreePosition(p)
	if err != nil {
		return false, err
	}
	if len(t.AllPlayers()) == 2 {
		return true, t.startFirstGame()
	}
	return false, nil
}

func (t *Table) MultiSetTournamentInfo(p *Player, place int) {
	t.setPlayerTournamentInfo(p, place)
}

func (t *Table) FirstAvailablePosition() int {
	for _, seat := range t.Seats {
		if seat.IsFree() {
			return seat.Position
		}
	}
	return -1
}