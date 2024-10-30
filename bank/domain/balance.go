package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Balance struct {
	UserID          string `bson:"_id"`
	Username        string
	Picture         string
	Chips           int64
	Coins 			int64
	CreatedAt       int64 // millis
	UpdatedAt       int64 // millis
	LastOperationID primitive.ObjectID
	NonRatable      bool

	Version int64
}

type ProfileInfo struct {
	UserID   string
	Username string
	Picture  string
}

func NewBalance(info ProfileInfo) *Balance {
	createAt := time.Now().UnixNano() / 1e6
	return &Balance{
		UserID: info.UserID,
		Username: info.Username,
		Picture: info.Picture,
		CreatedAt: createAt,
		UpdatedAt: createAt,
	}
}

func (b *Balance) SetBalanceCalc(chips, coins int64, lastOperationID primitive.ObjectID) {
	b.Chips = chips
	b.Coins = coins
	b.UpdatedAt = time.Now().UnixNano() / 1e6
	b.LastOperationID = lastOperationID
}

func (b *Balance) UpdateProfile(username, picture string) {
	b.Username = username
	b.Picture = picture
}

func (b *Balance) HandleOperation(op *Operation) {
	switch op.Type {
	case Deposit:
		b.Chips += op.Chips
	case Withdraw:
		b.Chips -= op.Chips
	case DepositCoins:
		b.Coins += op.Coins
	case WithdrawCoins:
		b.Coins -= op.Coins
	}
	b.UpdatedAt = time.Now().UnixNano() / 1e6
	b.LastOperationID = op.ID
}

func (b *Balance) Deposit(opID primitive.ObjectID, chips int64) {
	b.Chips += chips
	b.UpdatedAt = time.Now().UnixNano() / 1e6
	b.LastOperationID = opID
}

// Can withdraw so that chips balance becomes negative
func (b *Balance) Withdraw(opID primitive.ObjectID, chips int64) {
	b.Chips -= chips
	b.UpdatedAt = time.Now().UnixNano() / 1e6
	b.LastOperationID = opID
}

func (b *Balance) IsEnoughChips(chips int64) bool {
	return b.Chips >= chips
}

func (b *Balance) IsEnoughCoins(coins int64) bool {
	return b.Coins >= coins
}
