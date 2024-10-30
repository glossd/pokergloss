package domain

var ErrNotAvailableInSitngo = E("no available in Sit&Go game type")

func NewTableSitAndGo(params NewTableParams, seats []*Seat, attrs TournamentAttributes, autoStart bool) (*Table, error) {
	t, err := NewTable(params)
	if err != nil {
		return nil, err
	}

	t.Type = SitngoType
	t.setSeatsForTournament(seats)
	t.TournamentAttributes = attrs

	if autoStart {
		err = t.startFirstGame()
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

type TournamentWinner struct {
	*Player
	Place int
	Prize int64
}

func (t *Table) setSeatsForTournament(seats []*Seat) {
	for _, s := range seats {
		t.Seats[s.Position] = s
	}
}

func (t *Table) leaveSitngo(p *Player) {
	if t.IsSurvival {
		t.nullifyPlayer(p)
		return
	}

	place := len(t.AllPlayers())
	prize := t.setPlayerTournamentInfo(p, place)
	if prize > 0 {
		t.TournamentWinners = append(t.TournamentWinners, &TournamentWinner{Player: p, Prize: prize, Place: place})
	}

	t.nullifyPlayer(p)
}

func (t *Table) nullifyBrokePlayers() {
	players := t.PlayersFilter(BrokePlayerFilter)
	for _, player := range players {
		t.nullifyPlayer(player)
	}
}
