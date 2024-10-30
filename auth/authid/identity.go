package authid

import (
	"firebase.google.com/go/auth"
	"fmt"
)

var ErrUsernameNotSet = fmt.Errorf("username is not set")

type Identity struct {
	UserId string `json:"userId"`
	Username string `json:"username"`
	Picture string `json:"picture"`
}

func (id Identity) String() string {
	return fmt.Sprintf("{userId:%s, username:%s}", id.UserId, id.Username)
}

func FromRecord(record *auth.UserRecord) (*Identity, error) {
	if record == nil {
		return nil, fmt.Errorf("record is nil")
	}

	var username string
	usernameV := record.CustomClaims["username"]
	if v, ok := usernameV.(string); ok {
		username = v
	} else {
		return nil, ErrUsernameNotSet
	}

	return &Identity{
		UserId:  record.UID,
		Username: username,
		Picture:  record.PhotoURL,
	}, nil
}
