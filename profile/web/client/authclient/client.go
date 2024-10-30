package authclient

import (
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/profile/conf"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var ErrSetUsername = errors.New("couldn't set your username in database")

var AuthClient IAuthClient = &FirebaseAuthClient{}

func Init() {
	if goconf.IsLocal() {
		AuthClient = &MockAuthClient{}
	} else {
		conf.InitAuthClient()
		AuthClient = &FirebaseAuthClient{}
	}
}

type IAuthClient interface {
	CreateUser(ctx context.Context, username string, email string, password string) (*auth.UserRecord, error)
	GetCustomToken(ctx context.Context, uid string) (string, error)
	ChangeUsername(uid, username string, updateDisplayName bool) error
}

type FirebaseAuthClient struct{}

func (c FirebaseAuthClient) GetCustomToken(ctx context.Context, uid string) (string, error) {
	return conf.AuthClient.CustomToken(ctx, uid)
}

func (c FirebaseAuthClient) CreateUser(ctx context.Context, username string, email string, password string) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(username).
		EmailVerified(true)

	user, err := conf.AuthClient.CreateUser(ctx, params)
	if err != nil {
		log.Errorf("Couldn't create user: %s", err)
		return nil, err
	}

	err = c.ChangeUsername(user.UID, username, false)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (FirebaseAuthClient) ChangeUsername(uid, username string, updateDisplayName bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := conf.AuthClient.SetCustomUserClaims(ctx, uid, map[string]interface{}{"username": username})
	if err != nil {
		log.Errorf("Couldn't set username claim: %s", err)
		return ErrSetUsername
	}

	if updateDisplayName {
		_, err = conf.AuthClient.UpdateUser(ctx, uid, (&auth.UserToUpdate{}).DisplayName(username))
		if err != nil {
			log.Errorf("Display name wasn't updated for user, uid=%s, newUsername=%s : %s", uid, username, err)
		}
	}

	return nil
}

type MockAuthClient struct {
	UserID string
}

func (c MockAuthClient) GetCustomToken(ctx context.Context, uid string) (string, error) {
	return "customToken", nil
}

func (c MockAuthClient) CreateUser(ctx context.Context, username string, email string, password string) (*auth.UserRecord, error) {
	userID := primitive.NewObjectID().Hex()
	if c.UserID != "" {
		userID = c.UserID
	}
	return &auth.UserRecord{UserInfo: &auth.UserInfo{UID: userID}}, nil
}

func (MockAuthClient) ChangeUsername(uid, username string, updateDisplayName bool) error {
	return nil
}
