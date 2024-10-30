package domain

type Score struct {
	UserID string `bson:"_id"`
	// Reached.
	Level int
}
