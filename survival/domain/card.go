package domain

type Card struct {
	UserID string `bson:"_id"`
	// > 0, MongoDB validation.
	Tickets int64
}
