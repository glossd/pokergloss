package service

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/profile/conf"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/web/client/gcs"
	"github.com/glossd/pokergloss/profile/web/client/myfirestore"
	"github.com/glossd/pokergloss/profile/web/mq"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
)

func UpdateUserAvatar(ctx context.Context, file *multipart.FileHeader, iden authid.Identity) (string, error) {
	avatarURL := getAvatarURL(iden)
	var publicURL string
	var err error
	if goconf.IsProd() {
		publicURL, err = gcs.UploadFile(ctx, file, goconf.Props.AvatarBucket, avatarURL)
		if err != nil {
			log.Errorf("Couldn't update GCS file: %s", err)
			return "", err
		}
	} else {
		publicURL, err = myfirestore.UploadFile(ctx, file, goconf.Props.AvatarBucket, avatarURL)
		if err != nil {
			log.Errorf("Couldn't update firestore file: %s", err)
			return "", err
		}
	}

	_, err = conf.AuthClient.UpdateUser(ctx, iden.UserId, (&auth.UserToUpdate{}).PhotoURL(publicURL))
	if err != nil {
		log.Errorf("Couldn't update user: %s", err)
		return "", err
	}

	err = db.UpdatePicture(iden.Username, publicURL)
	if err != nil {
		log.Errorf("Couldn't update user picture in mongo: %s, uid=%s, url=%s", err, iden.UserId, publicURL)
	}

	mq.PublishProfileUpdateFields(iden.UserId, iden.Username, publicURL)

	return publicURL, nil
}

func getAvatarURL(iden authid.Identity) string {
	hash := uniuri.NewLen(4)
	return fmt.Sprintf("%s-%s", iden.UserId, hash)
}
