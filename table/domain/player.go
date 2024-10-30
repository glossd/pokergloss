package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/goconf/timeutil"
	log "github.com/sirupsen/logrus"
)

var ErrNotYourTurn = E("it's not you turn to act")
var ErrNotEnoughChips = E("you don't have enough chips on the table")

type Player struct {
	authid.Identity
	// table position from Seat.Position
	Position int

	Stack          int64
	StartGameStack int64 // to calculate winners in tournaments
	BuyInStack     int64
	Status         PlayerStatus
	Cards          *HoleCards
	Blind          Blind
	// Action made by a player in current round. round action should be set to nil on each new round
	LastRoundAction ActionType
	// sets to nil on each new round
	TotalRoundBet int64

	LastGameAction ActionType
	LastGameBet    int64

	Intent          *Intent
	isIntentRemoved bool
	isIntentChanged bool

	// Computes in Table.ComputeWinners. The less rank the stronger hand.
	HandRank       int32
	HandRankString string

	// Player don't leave the table right away if table is playing.
	// If a player is leaving, the he will be automatically folded and only then left.
	IsLeaving bool

	ChipsToAddOnReset int64

	SatOutAt int64

	AutoConfig PlayerAutoConfig

	ShowDownAction ShowDownActionType

	// Enrichment
	Level                   *int64
	BankBalance             *int64
	BankRank                *int64
	MarketItemID            string
	MarketItemCoinsDayPrice int64

	// Updates at player making action
	UpdatedAt int64

	tournamentInfo PlayerTournamentInfo

	isStandFolded         bool
	potIndex              int
	isMultiMoving         bool
	multiPreviousPosition int
	wonChips              int64
}

type PlayerStatus string

const (
	PlayerReservedSeat PlayerStatus = "sitting"
	PlayerReady        PlayerStatus = "ready"
	PlayerPlaying      PlayerStatus = "playing"
	PlayerSittingOut   PlayerStatus = "sittingOut"
)

type PlayerTournamentInfo struct {
	Place       int
	Prize       int64
	IsLast      bool
	marketPrize *MarketPrize
}

// Creates player in reserved seat state
func NewPlayer(iden authid.Identity, position int) *Player {
	return &Player{
		Identity:   iden,
		Status:     PlayerReservedSeat,
		Position:   position,
		AutoConfig: NewPlayerAutoConfig(iden.UserId),
	}
}

// Player just reserved a seat and put stack on a table
func (p *Player) setInitStack(stack int64) {
	p.Stack = stack
	p.StartGameStack = stack
	p.BuyInStack = stack
	p.Status = PlayerReady
}

func (p *Player) setHoleCards(f Card, s Card) {
	p.Cards = &HoleCards{
		First:  f,
		Second: s,
	}
}

func (p *Player) IsSittingOut() bool {
	return p != nil && p.Status == PlayerSittingOut
}

func (p *Player) setToSittingOut() {
	p.Status = PlayerSittingOut
	p.SatOutAt = timeutil.NowAdd(0)
}

func (p *Player) setAutoMuck(autoMuck bool) {
	p.AutoConfig.Muck = autoMuck
}

func (p *Player) setAutoTopUp(autoTopUp bool) {
	p.AutoConfig.TopUp = autoTopUp
}

func (p *Player) setAutoReBuy(autoReBuy bool) {
	p.AutoConfig.ReBuy = autoReBuy
}

func (p *Player) setShowDownAction(action ShowDownActionType) {
	p.ShowDownAction = action
}
func (p *Player) HasCards() bool {
	if p.Cards != nil {
		return true
	}
	return false
}

func (p *Player) makeChipsFreeAction(a Action) {
	p.LastRoundAction = a.Type
	p.LastGameAction = a.Type
	if a.IsChipFree() {
		p.LastGameBet = p.TotalRoundBet
	}
	p.UpdatedAt = timeutil.Now()
}

func (p *Player) newRoundReset() {
	p.LastRoundAction = ""
	p.TotalRoundBet = 0
	p.Intent = nil
	p.isIntentRemoved = false
	p.isIntentChanged = false
}

func (p *Player) gameReset() {
	p.newRoundReset()
	p.Cards = nil
	p.LastGameAction = ""
	p.LastGameBet = 0
	p.ShowDownAction = ""
	p.Blind = ""
	p.Stack += p.ChipsToAddOnReset
	p.StartGameStack = p.Stack
	p.ChipsToAddOnReset = 0
	p.HandRank = 0
	p.HandRankString = ""
}

func (p *Player) IsBroke() bool {
	return p.Status == PlayerSittingOut && p.Stack == 0
}

func (p *Player) addChipsToStack(chips int64) {
	p.Stack += chips
}

func (p *Player) addChipsOnGameStart(chips int64) {
	p.ChipsToAddOnReset += chips
}

func (p *Player) moveChipsFromBetToStack(chips int64) error {
	if p.TotalRoundBet < chips {
		return E("Player.moveChipsFromBetToStack not enough chips")
	}
	p.Stack += chips
	p.TotalRoundBet -= chips
	p.LastGameBet -= chips
	if p.LastGameAction == AllIn {
		p.LastGameAction = Call
		p.LastRoundAction = Call
	}
	return nil
}

func (p *Player) betChips(chips int64) error {
	if p.Stack < chips {
		log.Errorf("User tried to withdraw more chips than he has on table, identity=%s", p.Identity)
		return ErrNotEnoughChips
	}
	p.Stack -= chips
	p.TotalRoundBet += chips
	p.LastGameBet = p.TotalRoundBet
	return nil
}

// Returns bet chips and is bet AllIn type
func (p *Player) betChipsForBlind(blindForceBet int64) (int64, bool) {
	if p.Stack <= blindForceBet {
		betAllIn := p.Stack
		p.TotalRoundBet = betAllIn
		p.LastGameBet = betAllIn
		p.Stack = 0
		p.LastRoundAction = AllIn
		p.LastGameAction = AllIn
		return betAllIn, true
	} else {
		_ = p.betChips(blindForceBet)
		return blindForceBet, false
	}
}

func (p *Player) addWinningsToStack(chips int64) {
	p.Stack += chips
	p.wonChips += chips
}

func (p *Player) GetWonChips() int64 {
	return p.wonChips
}

func (p *Player) HasIntent() bool {
	return p.Intent != nil
}

func (p *Player) setIntent(intent Intent) {
	p.Intent = &intent
}

func (p *Player) removeIntent() {
	p.Intent = nil
	p.isIntentRemoved = true
}

func (p *Player) updateIntent(newMaxRoundBet, oldMaxRoundBet int64) {
	result := p.Intent.upgradeIntent(newMaxRoundBet, oldMaxRoundBet, p.Stack)
	switch result {
	case DeleteIntentUpgrade:
		p.removeIntent()
	case ChangeIntentUpgrade:
		p.isIntentChanged = true
	}
}

func (p *Player) IsIntentRemoved() bool {
	return p.isIntentRemoved
}

func (p *Player) IsIntentChanged() bool {
	return p.isIntentChanged
}

func (p *Player) GetIntentActionAndDelete() Action {
	intent := p.Intent
	p.removeIntent()
	return intent.Action
}

func (p *Player) AllChips() int64 {
	return p.Stack + p.ChipsToAddOnReset
}

func (p *Player) IsPlaying() bool {
	return p != nil && p.Status == PlayerPlaying
}

func (p *Player) IsReady() bool {
	return p != nil && p.Status == PlayerReady
}

func (p *Player) IsGaming() bool {
	return p.IsPlaying() && p.LastGameAction != Fold
}

// Player.Status can be PlayerSittingOut in
// tournament, and the action is fold (so called auto-fold)
func (p *Player) IsRoundFolded() bool {
	return p != nil && p.LastRoundAction == Fold
}

func (p *Player) IsDecidable() bool {
	return p.IsGaming() && p.LastGameAction != AllIn
}

func (p *Player) IsActive() bool {
	return p.IsPlaying() || p.Status == PlayerReady
}

func (p *Player) IsAvailableForNextGame() bool {
	return p.IsActive() && p.Stack > 0 && !p.IsLeaving
}

func (p *Player) IsAllIn() bool {
	return p.IsPlaying() && p.LastGameAction == AllIn
}

func (p *Player) IsMucked() bool {
	return p.ShowDownAction == Muck
}

func (p *Player) GetTournamentInfo() PlayerTournamentInfo {
	return p.tournamentInfo
}

func (p *Player) GetMultiPreviousPosition() int {
	return p.multiPreviousPosition
}

func (p *Player) IsStandFolded() bool {
	return p.isStandFolded
}

func (p *Player) setStartGameStack() {
	p.StartGameStack = p.Stack + p.TotalRoundBet // stack + blind chips
}

func (p *Player) SetLevel(level int64) {
	p.Level = &level
}

func (p *Player) SetBankInfo(chips, rank int64) {
	p.BankBalance = &chips
	p.BankRank = &rank

}

func (p *Player) SetMarketItem(itemID string, coinsDayPrice int64) {
	p.MarketItemID = itemID
	p.MarketItemCoinsDayPrice = coinsDayPrice
}

func (pti *PlayerTournamentInfo) GetTournamentMarketPrize() *MarketPrize {
	return pti.marketPrize
}

func (p *Player) SetTournamentMarketPrize(mp *MarketPrize) {
	p.tournamentInfo.marketPrize = mp
}

func (p *Player) GetAutoConfig() *PlayerAutoConfig {
	return &p.AutoConfig
}

func (p *Player) IsTopUpable() bool {
	return p.Stack < p.BuyInStack && p.Stack != 0 && !p.IsSittingOut()
}

// Please be careful
func (p *Player) TopUp() {
	if p.IsTopUpable() {
		p.Stack = p.BuyInStack
	} else {
		log.Error("domain.Player#TopUp: failed to top up")
	}
}

func (p *Player) IsReBuyable() bool {
	return p.Stack == 0
}

// Please be careful
func (p *Player) ReBuy() {
	if p.IsReBuyable() {
		p.Stack = p.BuyInStack
	} else {
		log.Error("domain.Player#ReBuy: failed to rebuy")
	}
}

func (p *Player) additionalDecisionTimePercent() float64 {
	switch p.MarketItemID {
	case "hourglass":
		return 0.2
	}
	if p.IsVip() {
		return 0.1
	}
	return 0
}

func (p *Player) IsVip() bool {
	return p.MarketItemCoinsDayPrice > 0
}
