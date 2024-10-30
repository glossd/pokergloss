package model

type BalanceWithRank struct {
	Balance int64 `json:"balance"`
	Coins int64 `json:"coins"`
	Rank int64 `json:"rank"`
}
