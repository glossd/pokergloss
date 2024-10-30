package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

var (
	ErrEmptyName = E("name can't be empty")
	ErrNameSize  = E("name can't be more than 20 letters")
)

var SystemIdentity = authid.Identity{UserId: "System", Username: "system"}

type Table struct {
	ID              primitive.ObjectID `bson:"_id"`
	Name            string
	Size            int
	BigBlind        int64
	SmallBlind      int64
	DecisionTimeout time.Duration
	CreatedAt       int64
	// userID
	CreateBy string

	Type TableType
	BettingLimit
	FixedLimitBet int64

	// Seats are created with creation of table meaning a seat can't be nil.
	Seats  []*Seat
	Status TableStatus

	RoundPot int64
	// The total pot includes the bet(s) made in the current round
	TotalPot        int64
	Pots            Pots
	lastTakenPotIdx int

	// Never nil
	CommunityCards *CommunityCards

	Winners []Winner

	// 1 would take 100% from winnings, 0 would be rake free
	RakePercent float64

	DecidingPosition int
	// It's also use for game end timeout.
	// millis
	DecisionTimeoutAt int64

	LastAggressorPosition int

	WaitingAt int64

	TournamentAttributes

	MultiAttrs

	// Optimistic locking.
	// Basically all table actions are serial.
	// Increased on playerAction, decisionTimeout and gameStart
	GameFlowVersion int64

	// is set on db operation, you can't use it inside domain
	// created for fast filtering of empty and full tables
	PlayersCount int

	// tells whether table can't be deleted
	IsPersistent bool

	IsPrivate bool

	IsSurvival    bool
	SurvivalLevel int64
	ThemeID       ThemeID

	// fields below don't need to be saved to db
	nullifiedLeavingPlayers   []*Player
	showedDownPlayers         []*Player
	brokePlayers              []*Player
	wasNewRound               bool
	previousDecidingPosition  int
	wasTimeout                bool
	madeActionPlayerPositions []*Player
	stackOverflowPlayer       *Player
	rake                      Rake
	// case when two player left and one goes allIn with blind
	isAutoGameEnd bool
}

type TableStatus string

const (
	WaitingTable  TableStatus = "waiting"
	PlayingTable  TableStatus = "playing"
	GameEndTable  TableStatus = "gameEnd"
	ShowdownTable TableStatus = "showDown"
)

const DefaultDecisionTimeout = 10 * time.Second

// made it public for testing, don't ever change it in production code
var MaxTableSize = 20

type NewTableParams struct {
	Name            string
	Size            int
	BigBlind        int64
	DecisionTimeout time.Duration
	BettingLimit
	IsPrivate bool
	authid.Identity
}

// If decisionTimeout equals 0 it sets to DefaultDecisionTimeout
func NewTable(params NewTableParams) (*Table, error) {
	err := ValidateAndEnrichParams(&params)
	if err != nil {
		return nil, err
	}

	seats := make([]*Seat, params.Size)
	for i := 0; i < params.Size; i++ {
		seats[i] = newSeat(i)
	}
	newID := primitive.NewObjectID()
	if strings.TrimSpace(params.UserId) == "" {
		log.Warnf("NewTable wasn't identified, tableID=%s", newID.Hex())
	}
	return &Table{
		ID:                    newID,
		Name:                  params.Name,
		Size:                  params.Size,
		BigBlind:              params.BigBlind,
		SmallBlind:            params.BigBlind / 2,
		DecisionTimeout:       params.DecisionTimeout,
		BettingLimit:          params.BettingLimit,
		IsPrivate:             params.IsPrivate,
		Type:                  CashType,
		Seats:                 seats,
		CreatedAt:             timeutil.Now(),
		CreateBy:              params.UserId,
		Status:                WaitingTable,
		WaitingAt:             timeutil.Now(),
		CommunityCards:        &CommunityCards{},
		DecidingPosition:      -1,
		LastAggressorPosition: -1,
		GameFlowVersion:       0,
		Pots:                  initPots(),
		RakePercent:           GetRakePercent(),
	}, nil
}

func ValidateAndEnrichParams(params *NewTableParams) error {
	normalizedName := strings.TrimSpace(params.Name)
	if normalizedName == "" {
		return ErrEmptyName
	}
	if len(normalizedName) > 21 {
		return ErrNameSize
	}
	params.Name = normalizedName

	if params.Size < 2 {
		return E("size must be more than 1")
	}
	if params.Size > MaxTableSize {
		return E("size must be less than %d", MaxTableSize+1)
	}

	if params.BigBlind < 2 {
		return E("table is not allowed to have big blind less than 2 chips")
	}

	if params.DecisionTimeout == 0 {
		params.DecisionTimeout = DefaultDecisionTimeout
	}

	if params.DecisionTimeout < conf.Props.MinDecisionTimeout || params.DecisionTimeout > conf.Props.MaxDecisionTimeout {
		log.Warnf("User send not allowed decision timeout = %d seconds", params.DecisionTimeout)
		return E("timeout must be less then %s and more than %s seconds",
			conf.Props.MaxDecisionTimeout, conf.Props.MinDecisionTimeout)
	}

	if params.BettingLimit == "" {
		params.BettingLimit = NL
	}

	return nil
}

func (t *Table) setToWaiting() {
	t.Status = WaitingTable
	t.WaitingAt = timeutil.Now()
	if t.IsCashType() {
		for _, player := range t.PlayingPlayersByGameType() {
			player.Status = PlayerReady
		}
	}
	log.Infof("Table has stopped id=%s", t.ID.Hex())
}

func (t *Table) setToGameEnd() {
	t.Status = GameEndTable
}

func (t *Table) setToShowDown() {
	t.Status = ShowdownTable
}

func (t *Table) IsInPlay() bool {
	return t.IsPlaying() || t.IsShowDown()
}

func (t *Table) IsActive() bool {
	return t.IsInPlay() || t.Status == GameEndTable
}

func (t *Table) IsWaiting() bool {
	return t.Status == WaitingTable
}

func (t *Table) IsPlaying() bool {
	return t.Status == PlayingTable
}

func (t *Table) IsGameEnd() bool {
	return t.Status == GameEndTable
}

func (t *Table) IsShowDown() bool {
	return t.Status == ShowdownTable
}

func (t *Table) IsCashType() bool {
	return t.Type == CashType
}

func (t *Table) IsSitngoType() bool {
	return t.Type == SitngoType
}

func (t *Table) IsTournament() bool {
	return t.Type == SitngoType || t.Type == MultiType
}

func (t *Table) IsMultiType() bool {
	return t.Type == MultiType
}

func (t *Table) IsRingType() bool {
	return t.Type == CashType
}

func (t *Table) MaxBuyInStack() int64 {
	return MaxBuyIn(t.BigBlind)
}

func MaxBuyIn(bigBlind int64) int64 {
	return bigBlind * 200
}

func (t *Table) MinBuyInStack() int64 {
	return MinBuyIn(t.BigBlind)
}

func MinBuyIn(bigBlind int64) int64 {
	return bigBlind * 50
}

func (t *Table) MaxRoundBet() int64 {
	var max int64
	for _, player := range t.PlayersFilter(t.GameTypePlayerFilter()) {
		if max < player.TotalRoundBet {
			max = player.TotalRoundBet
		}
	}
	return max
}

func (t *Table) NullifiedLeavingPlayerPositions() []int {
	var positions []int
	for _, player := range t.nullifiedLeavingPlayers {
		if player != nil {
			positions = append(positions, player.Position)
		}
	}
	return positions
}

func (t *Table) NullifiedLeavingPlayers() []*Player {
	return t.nullifiedLeavingPlayers
}

func (t *Table) BrokePlayers() []*Player {
	return t.brokePlayers
}

func (t *Table) IsAutoGameEnd() bool {
	return t.isAutoGameEnd
}

func (t *Table) MultiSortedNullifiedPlayers() []*Player {
	SortPlayersByStartGameStack(t.nullifiedLeavingPlayers)
	if t.MultiAttrs.IsLast {
		for i, p := range t.nullifiedLeavingPlayers {
			if p.tournamentInfo.IsLast && i != len(t.nullifiedLeavingPlayers)-1 {
				newSortedPlayers := append(t.nullifiedLeavingPlayers[:i], t.nullifiedLeavingPlayers[i+1:]...)
				t.nullifiedLeavingPlayers = append(newSortedPlayers, p)
			}
		}
	}
	return t.nullifiedLeavingPlayers
}
func (t *Table) ShowedDownPlayers() []*Player {
	return t.showedDownPlayers
}

func (t *Table) nullifyStandingPlayers() {
	players := t.PlayersFilter(func(p *Player) bool { return p.IsLeaving })
	for _, player := range players {
		t.nullifyPlayer(player)
	}
}

func (t *Table) nullifyPlayer(p *Player) {
	t.nullifiedLeavingPlayers = append(t.nullifiedLeavingPlayers, p)
	t.removePlayer(p)
}

func (t *Table) MadeActionPlayerPositions() []*Player {
	return t.madeActionPlayerPositions
}

func (t *Table) lastAggressorOrFirstRoundPosition() (int, error) {
	if t.LastAggressorPosition >= 0 {
		return t.LastAggressorPosition, nil
	}

	var firstRoundP *Player
	var err error
	if t.IsPreFlop() {
		firstRoundP, err = t.nextPlayer(t.BigBlindPosition(), t.nextPlayerToDecideFilter())
	} else {
		firstRoundP, err = t.nextPlayer(t.DealerPosition(), t.nextPlayerToDecideFilter())
	}
	if err != nil {
		p, onlyOne := t.isOneGamingLeftPlayer()
		if onlyOne {
			return p.Position, nil
		}
		return 0, err
	}

	return firstRoundP.Position, nil
}

func (t *Table) setToDecidingNoTimeout(p *Player) {
	t.previousDecidingPosition = t.DecidingPosition
	t.DecidingPosition = p.Position
}

func (t *Table) IsPositionExists(pos int) bool {
	err := t.validatePosition(pos)
	return err == nil
}

func (t *Table) IsSeatTaken(pos int) bool {
	return t.IsSeatFree(pos)
}

func (t *Table) AllCards() []Card {
	var allCards []Card
	for _, seat := range t.Seats {
		if seat.IsTaken() {
			holeCards := seat.GetPlayer().Cards
			if holeCards != nil {
				allCards = append(allCards, holeCards.Get()...)
			}
		}
	}
	allCards = append(allCards, t.CommunityCards.AvailableCards()...)
	return allCards
}

func (t *Table) IsDeciding(p *Player) bool {
	if p == nil {
		return false
	}
	return t.DecidingPosition == p.Position
}

func (t *Table) NullifySittingOutPlayer(position int) error {
	p, err := t.GetPlayer(position)
	if err != nil {
		return err
	}

	if !p.IsSittingOut() {
		return E("domain.NullifySittingOutPlayer: can't nullify non sitting out player")
	}

	t.nullifyPlayer(p)
	return nil
}

// Nillable.
// This function exists for the sake of showing
// the stack before the end result.
func (t *Table) BuildStackOverflowPlayer() *Player {
	if t.stackOverflowPlayer == nil {
		return nil
	}
	if t.IsGameEnd() {
		sofp := t.copyStackOverflowPlayer()
		sofp.Stack = t.getStackBeforeWinning(&sofp)
		return &sofp
	}
	return t.stackOverflowPlayer
}

func (t *Table) GetPlayerStackBeforeResult(p *Player) int64 {
	if p == nil {
		return 0
	}
	if t.IsGameEnd() {
		return t.getStackBeforeWinning(p)
	}
	return p.Stack
}

func (t *Table) copyStackOverflowPlayer() Player {
	if t.stackOverflowPlayer == nil {
		log.Fatalf("You can't use copyStackOverflowPlayer on nil")
	}
	return *t.stackOverflowPlayer
}

func (t *Table) getStackBeforeWinning(p *Player) int64 {
	for _, winner := range t.Winners {
		if winner.Position == p.Position {
			rake := t.buildRake()
			return p.Stack - winner.Chips + rake.Of(p.Position)
		}
	}
	return p.Stack
}

func (t *Table) isWinner(pos int) bool {
	return t.Pots.isWinnerOfPot(pos)
}

func (t *Table) SetPlayerAutoConfig(p *Player, config PlayerAutoConfig) {
	if t.IsSitngoType() || t.IsMultiType() {
		config.TopUp = false
		config.ReBuy = false
	}
	p.AutoConfig = config
}

func (t *Table) sitPlayerToReady(pos int, iden authid.Identity, stack int64) {
	t.Seats[pos].addPlayer(iden)
	t.Seats[pos].Player.setInitStack(stack)
}
