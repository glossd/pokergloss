package domain

import "github.com/glossd/pokergloss/auth/authid"

type EntrySitAndGo struct {
	authid.Identity
	Position int
}
