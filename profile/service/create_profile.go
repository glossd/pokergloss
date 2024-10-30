package service

import (
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/domain"
	"github.com/glossd/pokergloss/profile/web/client/authclient"
	"github.com/glossd/pokergloss/profile/web/mq"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var usernameRegexp, _ = regexp.Compile(`^[A-Za-z0-9_]+$`)

var ErrUsernameTaken = errors.New("username is taken")
var ErrEmptyUsername = errors.New("username can't be empty")
var ErrSpaceUsername = errors.New("username can't have empty spaces")
var ErrMaxLengthUsername = errors.New("username must be less than 20 letters")
var ErrMinLengthUsername = errors.New("username must be more than 2 letters")
var ErrUsernameWrong = errors.New(`username must contain latin letters, digits and "_" character`)

func CreateUser(ctx context.Context, username, email, password string, info domain.TechInfo) (*auth.UserRecord, error) {
	err := validateUsername(username)
	if err != nil {
		return nil, err
	}

	newProfile := domain.NewProfile(username, "", info)
	err = db.SaveProfile(ctx, newProfile)
	if err != nil {
		if strings.HasPrefix(err.Error(), "multiple write errors: [{write errors: [{E11000 duplicate key") {
			return nil, ErrUsernameTaken
		}
		return nil, err
	}

	user, err := authclient.AuthClient.CreateUser(ctx, username, email, password)
	if err != nil {
		return nil, err
	}

	newProfile.UserID = user.UID
	err = db.UpsertProfileNoCtx(newProfile)
	if err != nil {
		log.Errorf("Couldn't set userID for username doc, username=%s, userID=%s: %s", username, user.UID, err)
	}

	mq.PublishCreated(newProfile)
	mq.PublishProfileUpdate(newProfile)

	return user, err
}

func validateUsername(username string) error {
	if username == "" {
		return ErrEmptyUsername
	}

	if strings.ContainsRune(username, ' ') {
		return ErrSpaceUsername
	}

	if len(username) >= 20 {
		return ErrMaxLengthUsername
	}

	if len(username) < 3 {
		return ErrMinLengthUsername
	}

	if strings.HasPrefix(username, "_") {
		return errors.New(`username can't start with "_" underscore`)
	}

	if strings.HasSuffix(username, "_") {
		return errors.New(`username can't end with "_" underscore`)
	}

	if strings.Contains(username, "__") {
		return errors.New(`username can't have multiple "_" underscores in a row `)
	}

	if _, ok := forbiddenUsernames[username]; ok {
		return ErrUsernameTaken
	}

	if !usernameRegexp.MatchString(username) {
		return ErrUsernameWrong
	}
	return nil
}

var forbiddenUsernames = map[string]struct{}{
	"Demon":      {},
	"Archdemon":  {},
	"ArchDemon":  {},
	"Raziel":     {},
	"Fog_Spirit": {},
	"Torturer":   {},
}

func GetCustomToken(ctx context.Context, uid string) (string, error) {
	return authclient.AuthClient.GetCustomToken(ctx, uid)
}
