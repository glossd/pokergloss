package domain

type PlayerAutoConfig struct {
	UserID string `bson:"_id"`
	Muck bool
	TopUp bool
	ReBuy bool
}

func NewPlayerAutoConfig(userID string) PlayerAutoConfig {
	return PlayerAutoConfig{
		UserID: userID,
		Muck: true,
		TopUp: false,
		ReBuy: false,
	}
}
