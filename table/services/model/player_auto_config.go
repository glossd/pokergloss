package model

import "github.com/glossd/pokergloss/table/domain"

type PlayerAutoConfig struct {
	Muck  bool `json:"autoMuck"`
	TopUp bool `json:"autoTopUp"`
	ReBuy bool `json:"autoRebuy"`
}

func ToPlayerAutoConfig(pac *domain.PlayerAutoConfig) *PlayerAutoConfig {
	return &PlayerAutoConfig{
		Muck:  pac.Muck,
		TopUp: pac.TopUp,
		ReBuy: pac.ReBuy,
	}
}
