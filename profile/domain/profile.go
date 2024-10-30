package domain

import "time"

type Profile struct {
	Username string `bson:"_id"`
	UserID string
	Picture string
	OldUsernames []OldUsername
	CreatedAt int64
	UpdatedAt int64
	TechInfo
}

type OldUsername struct {
	Username string
	ChangedAt int64
}

type TechInfo struct {
	IP string
	Lang string
	Browser string
	OS string
}

func NewProfile(username, userID string, info TechInfo) *Profile {
	now := time.Now().Unix()
	return &Profile{Username: username, UserID: userID, TechInfo: info, CreatedAt: now, UpdatedAt: now}
}

func NewFromOld(newUsername string, old *Profile) *Profile {
	now := time.Now().Unix()
	oldUsernames := append(old.OldUsernames, OldUsername{ChangedAt: time.Now().Unix(), Username: old.Username})
	return &Profile{Username: newUsername, UserID: old.UserID, Picture: old.Picture, OldUsernames: oldUsernames, TechInfo: old.TechInfo, CreatedAt: old.CreatedAt, UpdatedAt: now}
}
