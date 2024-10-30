package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var ErrNegativeChips = E("number of chips must be positive")
var ErrNegativeCoins = E("number of coins must be positive")
var ErrNotEnoughChips = E("you don't have enough chips")
var ErrNotEnoughCoins = E("you don't have enough coins")

type Operation struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Type   OperationType
	Reason Reason
	// > 0
	Chips int64
	// > 0
	Coins int64
	UserID string // todo add index?
	Description string
	CreatedAt int64
}

type OperationType string

const (
	Deposit  OperationType = "deposit"
	Withdraw OperationType = "withdraw"
	DepositCoins OperationType = "depositCoins"
	WithdrawCoins OperationType = "withdrawCoins"
)

type Reason string

const (
	// deposit.
	Bonus       Reason = "bonus"
	NewLevel    Reason = "newLevel"
	Achievement Reason = "achievement"
	Assignment Reason = "assignment"
	Survival Reason = "survival"

	// withdraw.
	Game        Reason = "game"
	Market Reason = "market"
	InactionFee Reason = "inactionFee"

	Admin       Reason = "admin"
)

func NewDeposit(r Reason, chips int64, userID, description string) (*Operation, error){
	if chips <= 0 {
		return nil, ErrNegativeChips
	}
	return &Operation{
		Type:        Deposit,
		Reason:      r,
		Chips:       chips,
		UserID:      userID,
		Description: description,
		CreatedAt: now(),
	}, nil
}

func NewWithdraw(r Reason, chips int64, userID, description string) (*Operation, error) {
	if chips <= 0 {
		return nil, ErrNegativeChips
	}
	return &Operation{
		Type:        Withdraw,
		Reason:      r,
		Chips:       chips,
		UserID:      userID,
		Description: description,
		CreatedAt: now(),
	}, nil
}

func NewDepositCoins(r Reason, coins int64, userID, description string) (*Operation, error){
	if coins <= 0 {
		return nil, ErrNegativeCoins
	}
	return &Operation{
		Type:        DepositCoins,
		Reason: r,
		Coins:       coins,
		UserID:      userID,
		Description: description,
		CreatedAt: now(),
	}, nil
}

func NewWithdrawCoins(r Reason, coins int64, userID, description string) (*Operation, error){
	if coins <= 0 {
		return nil, ErrNegativeCoins
	}
	return &Operation{
		Type:        WithdrawCoins,
		Reason: r,
		Coins:       coins,
		UserID:      userID,
		Description: description,
		CreatedAt: now(),
	}, nil
}

func now() int64 {
	return time.Now().UnixNano() / 1e6
}
