package model

import (
	"github.com/glossd/pokergloss/profile/domain"
)

type Profile struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Picture  string `json:"picture"`
}

func ToProfiles(u []*domain.Profile) []*Profile {
	profiles := make([]*Profile, 0, len(u))
	for _, doc := range u {
		profiles = append(profiles, ToProfile(doc))
	}
	return profiles
}

func ToProfile(u *domain.Profile) *Profile {
	return &Profile{
		UserID:   u.UserID,
		Username: u.Username,
		Picture:  u.Picture,
	}
}
