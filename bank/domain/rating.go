package domain

type Rating struct {
	UserID string `bson:"_id"`
	Username string
	Picture string
	Chips int64
	Rank int64
	UpdatedAt int64
}
