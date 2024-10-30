package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/domain"
	"github.com/glossd/pokergloss/profile/web/client/authclient"
	"github.com/glossd/pokergloss/profile/web/mq"
	log "github.com/sirupsen/logrus"
	"strings"
)

func ChangeUsername(ctx context.Context, newUsername string, iden authid.Identity) error {
	err := validateUsername(newUsername)
	if err != nil {
		return err
	}

	profile, err := db.FindProfile(ctx, iden.Username)
	if err != nil {
		log.Errorf("Failed to find profile=%s : %s", newUsername, err)
		return err
	}

	newProfile := domain.NewFromOld(newUsername, profile)
	err = db.SaveProfile(ctx, newProfile)
	if err != nil {
		if strings.HasPrefix(err.Error(), "multiple write errors: [{write errors: [{E11000 duplicate key") {
			return ErrUsernameTaken
		}
		log.Errorf("Failed to update username profile=%s : %s", newUsername, err)
		return err
	}

	err = authclient.AuthClient.ChangeUsername(iden.UserId, newUsername, true)
	if err != nil {
		dErr := db.DeleteProfile(newUsername)
		if dErr != nil {
			log.Errorf("Couldn't delete new saved username=%s : %s", newUsername, err)
		}
		log.Errorf("Couldn't change username %s in GIDP: %s", newUsername, err)
		return err
	}

	err = db.DeleteProfile(iden.Username)
	if err != nil {
		log.Errorf("Failed to delete old username profile, oldUsername=%s", iden.Username)
	}

	mq.PublishProfileUpdate(newProfile)

	return nil
}
