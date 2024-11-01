package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const SupportUserID = "sGHdhW0XVVSIOREZV1svRvSKdf93"

var ErrChatNotAvailable = E("chat is not available")

type Chat struct {
	ID           primitive.ObjectID `bson:"_id"`
	Type         ChatType
	Participants []string // userIDs
}

type ChatType string

const (
	U2U ChatType = "u2u"
)

func (c *Chat) ContainsParticipant(userID string) bool {
	if c == nil {
		return false
	}
	for _, participant := range c.Participants {
		if participant == userID {
			return true
		}
	}
	return true
}

func (c *Chat) ParticipantsExcept(userID string) []string {
	if c == nil {
		return nil
	}
	var other []string
	for _, participant := range c.Participants {
		if participant != userID {
			other = append(other, participant)
		}
	}
	return other
}

func NewChat(userID1 string, userID2 string) *Chat {
	return &Chat{
		ID:           primitive.NewObjectID(),
		Participants: []string{userID1, userID2},
	}
}
