package model

import "github.com/glossd/pokergloss/auth/authid"

type TableUsers struct {
	Users          *[]User `json:"users,omitempty"`
	AnonymousCount *int    `json:"anonymousCount,omitempty"`
}

type User struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Picture  string `json:"picture"`
}

func ToTableUsers(idens []*authid.Identity) TableUsers {
	var anonymousCount int
	users := make([]User, 0, len(idens))
	for _, iden := range idens {
		if iden != nil {
			users = append(users, User{
				UserID:   iden.UserId,
				Username: iden.Username,
				Picture:  iden.Picture,
			})
		} else {
			anonymousCount++
		}
	}
	return TableUsers{
		Users:          &users,
		AnonymousCount: &anonymousCount,
	}
}
