package model

import "github.com/glossd/pokergloss/table/domain"

type UserWithStack struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Picture  string `json:"picture"`
	Stack    int64  `json:"stack"`
}

func ToUsersWithStack(ps []*domain.Player) []*UserWithStack {
	users := make([]*UserWithStack, 0, len(ps))
	for _, p := range ps {
		users = append(users, ToUserWithStack(p))
	}
	return users
}

func ToUserWithStack(p *domain.Player) *UserWithStack {
	return &UserWithStack{
		UserId:   p.UserId,
		Username: p.Username,
		Picture:  p.Picture,
		Stack:    p.Stack + p.TotalRoundBet,
	}
}
