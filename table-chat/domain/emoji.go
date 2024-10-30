package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Emoji string

const (
	Joy           Emoji = "joy"
	Wink          Emoji = "wink"
	Cry           Emoji = "cry"
	Rage          Emoji = "rage"
	Like          Emoji = "like"
	Scream        Emoji = "scream"
	Sunglasses    Emoji = "sunglasses"
	RaisedEyebrow Emoji = "raisedEyebrow"
)

var setOfEmoji = map[Emoji]struct{}{
	Joy: {}, Wink: {}, Cry: {}, Rage: {}, Like: {}, Scream: {}, Sunglasses: {}, RaisedEyebrow: {},
}

func IsValidEmoji(emoji string) bool {
	_, ok := setOfEmoji[Emoji(emoji)]
	return ok
}

type EmojiMessage struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	TableID string

	CreatedBy authid.Identity
	CreatedAt int64

	Emoji
}
