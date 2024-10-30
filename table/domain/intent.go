package domain

import log "github.com/sirupsen/logrus"

type IntentType string
const (

	FoldIntentType  IntentType = "fold"
	AllInIntentType IntentType = "all-in"

	// Someone bet before
	CallIntentType    IntentType = "call"
	CallFoldIntentType IntentType = "call-fold"
	CallAnyIntentType IntentType = "call-any"
	RaiseIntentType   IntentType = "raise"

	// Nobody bet before
	CheckFoldIntentType    IntentType = "check-fold"
	CheckIntentType        IntentType = "check"
	CheckCallAnyIntentType IntentType = "check-call-any"
	BetIntentType          IntentType = "bet"
)

type UpgradeResult string
const (
	DeleteIntentUpgrade UpgradeResult = "delete"
	SameIntentUpgrade   UpgradeResult = "same"
	ChangeIntentUpgrade UpgradeResult = "change"
)

type Intent struct {
	Type IntentType
	Action
}

var (
	FoldIntent = Intent{Type: FoldIntentType, Action: FoldAction}
	CheckIntent = Intent{Type: CheckIntentType, Action: CheckAction}
	CallIntent = Intent{Type: CallIntentType, Action: CallAction}
	CallFoldIntent = Intent{Type: CallFoldIntentType, Action: CallAction}
	CallAnyIntent = Intent{Type: CallAnyIntentType, Action: CallAction}
	AllInIntent = Intent{Type: AllInIntentType, Action: AllInAction}
)

func NewIntent(inType IntentType, chips int64) Intent {
	switch inType {
	case FoldIntentType:
		return FoldIntent
	case AllInIntentType:
		return AllInIntent
	case CallIntentType:
		return CallIntent
	case CallFoldIntentType:
		return CallFoldIntent
	case CallAnyIntentType:
		return CallAnyIntent
	case RaiseIntentType:
		return Intent{Type: inType, Action: RaiseAction(chips)}
	case CheckIntentType:
		return CheckIntent
	case CheckFoldIntentType, CheckCallAnyIntentType:
		return Intent{Type: inType, Action: CheckAction}
	case BetIntentType:
		return Intent{Type: inType, Action: BetAction(chips)}
	}

	log.Fatalf("Switch doesn't contain asked intentType: %s", inType)
	return Intent{}
}

var betBeforeIntents = map[IntentType]struct{}{
	FoldIntentType: {}, AllInIntentType: {}, CallIntentType: {}, CallFoldIntentType: {}, CallAnyIntentType: {}, RaiseIntentType: {},
}

var nobodyBetBeforeIntents = map[IntentType]struct{}{
	FoldIntentType: {}, AllInIntentType: {}, CheckFoldIntentType: {}, CheckIntentType: {}, CheckCallAnyIntentType: {}, BetIntentType: {},
}

func isBetBeforeIntent(intent IntentType) bool {
	_, ok := betBeforeIntents[intent]
	return ok
}

func isNobodyBetBeforeIntent(intent IntentType) bool {
	_, ok := nobodyBetBeforeIntents[intent]
	return ok
}

func (in *Intent) upgradeIntent(newMaxRoundBet, oldMaxRoundBet int64, stack int64) UpgradeResult {
	betDiff := newMaxRoundBet - oldMaxRoundBet
	if in.Type == FoldIntentType || in.Type == AllInIntentType {
		return SameIntentUpgrade
	}

	switch in.Type {
	case CallIntentType:
		// maybe try: in = nil
		return DeleteIntentUpgrade
	case CallFoldIntentType:
		in.upgradeToFold()
		return ChangeIntentUpgrade
	case CallAnyIntentType:
		if betDiff >= stack {
			in.upgradeToAllIn()
			return ChangeIntentUpgrade
		}
		return SameIntentUpgrade
	case RaiseIntentType:
		if in.Action.Chips > newMaxRoundBet*2 {
			return SameIntentUpgrade
		}
		return DeleteIntentUpgrade

	case CheckFoldIntentType:
		in.upgradeToFold()
		return ChangeIntentUpgrade
	case CheckIntentType:
		return DeleteIntentUpgrade
	case CheckCallAnyIntentType:
		if betDiff >= stack {
			in.upgradeToAllIn()
			return ChangeIntentUpgrade
		}
		in.upgradeToCallAny()
		return ChangeIntentUpgrade
	case BetIntentType:
		if in.Action.Chips > newMaxRoundBet*2 {
			in.upgradeToRaise()
			return ChangeIntentUpgrade
		}
		return DeleteIntentUpgrade
	}

	log.Errorf("upgradeIntent: switch doesn't have intent type %s", in.Type)
	return DeleteIntentUpgrade
}

func (in *Intent) IsAggressive() bool {
	return in.Type == BetIntentType || in.Type == RaiseIntentType || in.Type == AllInIntentType
}

func (in *Intent) upgradeToAllIn() {
	// try: in = AllInIntent
	in.Type = AllInIntentType
	in.Action = AllInAction
}

func (in *Intent) upgradeToFold() {
	// try: in = FoldIntent
	in.Type = FoldIntentType
	in.Action = FoldAction
}

func (in *Intent) upgradeToCallAny() {
	// try: in = CallAnyIntent
	in.Type = CallAnyIntentType
	in.Action = CallAction
}

func (in *Intent) upgradeToRaise() {
	in.Type = RaiseIntentType
	in.Action = RaiseAction(in.Action.Chips)
}
