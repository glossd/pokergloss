package e2e

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/stretchr/testify/assert"
"github.com/glossd/pokergloss/auth/authjwt"
"google.golang.org/api/option"
"log"
"strings"
"testing"
"time"
)

// Admin sdk doesn't have anything to retrieve id token
// to run the test successfully, please update the tokenStr
// or just don't bother :)

const (
	tokenStrWithoutPicture = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjQ5YWQ5YmM1ZThlNDQ3OTNhMjEwOWI1NmUzNjFhMjNiNDE4ODA4NzUiLCJ0eXAiOiJKV1QifQ.eyJ1c2VybmFtZSI6ImRlbmlzMiIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9wb2tlcmJsb3ciLCJhdWQiOiJwb2tlcmJsb3ciLCJhdXRoX3RpbWUiOjE1OTk4Mjg0MTksInVzZXJfaWQiOiJnWTVtZnIzOGRhUVhMdU94QzZtUHV6MExTNmEyIiwic3ViIjoiZ1k1bWZyMzhkYVFYTHVPeEM2bVB1ejBMUzZhMiIsImlhdCI6MTU5OTgyODQxOSwiZXhwIjoxNTk5ODMyMDE5LCJlbWFpbCI6ImRlbmlzZ2xvdG92LjE5MTFAbWFpbC5ydSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJlbWFpbCI6WyJkZW5pc2dsb3Rvdi4xOTExQG1haWwucnUiXX0sInNpZ25faW5fcHJvdmlkZXIiOiJwYXNzd29yZCJ9fQ.Pc8JlpqOlz3lCg-MpFGXQ-M_kAhvMk-WEw8I2170hFk5FCGQ-VqjDQ6rQZ78RFkWTtTBEDCYsj7kue2OmahFp-EWVQVbVT1PzwMx_Y7o8SIxUtxA-4K7p0vfRDUJFxY3vyQ-ebcnUl3kCaXS3aduy6DLXDmHNEKMi7OffAvZ_mR6CrJ8FEkcnS-eCzc7q-iWJuesGfnNebY-76_Dy_hwKaYuSe0AA3gjoVV_xsPnDvJpV5si-ZWFiH36bnjC4ngAl_DfOUKFkvDUnAiLElY4wYK5Xag1GhiAJmSjpelVNb8AXQfCpjnXrzHwP0VQT9yGAy0S3RuuTpMIh2Z8U5jKxA"
	tokenStr               = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjhmNDMyMDRhMTc5MTVlOGJlN2NjZDdjYjI2NGRmNmVhMzgzYzQ5YWIiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoicG9rZXIiLCJwaWN0dXJlIjoiaHR0cHM6Ly9zdG9yYWdlLmdvb2dsZWFwaXMuY29tL3Bva2VyYmxvdy1hdmF0YXJzL3NZbmRSa01GMlZNSFNHR2JIRXM1UWozczZqazItWGRweiIsInVzZXJuYW1lIjoicG9rZXIiLCJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vcG9rZXJibG93IiwiYXVkIjoicG9rZXJibG93IiwiYXV0aF90aW1lIjoxNjI1NTAzMjkxLCJ1c2VyX2lkIjoic1luZFJrTUYyVk1IU0dHYkhFczVRajNzNmprMiIsInN1YiI6InNZbmRSa01GMlZNSFNHR2JIRXM1UWozczZqazIiLCJpYXQiOjE2MjU2ODc0NzQsImV4cCI6MTYyNTY5MTA3NCwiZW1haWwiOiJkZW5pc2dsb3Rvdi4xOTExQG1haWwucnUiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJlbWFpbCI6WyJkZW5pc2dsb3Rvdi4xOTExQG1haWwucnUiXX0sInNpZ25faW5fcHJvdmlkZXIiOiJwYXNzd29yZCJ9fQ.QBIcsoXcYLlo0xLXNNX4pw_WG2gFNOutEvItVDa_CbuqmCnUDaFU_5PCO_eSvByNMznDjjTUH4yb8an1l0DPXlGoUU3u75EyXN68nRD6CVA7evPf-bsoLhMFojpa9OBJSkqzs8O-6kuuu4LxaiT1cWLg2j4p0eqy2FfaMNQ8IyWBrH7nUUw4URqU4jTpK8XYkj2Wnz3WZYworoR9KUXwnV4hSMcr4W0mJgBkMLxSjMTQ_Axmswexa_jyWZYXMhQAzsBN4A46tXGepCG3ILTh1lrHrVsbnOR2PxTjuDPYBDQwH3EpzvX2Juz4oNcUC-1umojwKQYLeQxnOyHF_Ice7g"
)

func TestVerifyToken(t *testing.T) {
	c := initFirebaseClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id, err := authjwt.VerifyToken(ctx, c, tokenStr)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ID token has expired at") {
			return
		}
	}
	assert.Nil(t, err)

	assert.Equal(t, "poker", id.Username)
	assert.Equal(t, "sYndRkMF2VMHSGGbHEs5Qj3s6jk2", id.UserId)
	assert.Equal(t, "https://storage.googleapis.com/pokerblow-avatars/sYndRkMF2VMHSGGbHEs5Qj3s6jk2-Xdpz", id.Picture)
	assert.True(t, id.EmailVerified)
	assert.Equal(t, "password", id.Provider)
}

func TestIsEmailVerified(t *testing.T) {
	c := initFirebaseClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	token, err := authjwt.VerifyToken(ctx, c, tokenStr)
	assert.Nil(t, err)
	assert.True(t, token.EmailVerified)
}

func TestVerifyTokenWithoutPicture(t *testing.T) {
	c := initFirebaseClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id, err := authjwt.VerifyToken(ctx, c, tokenStrWithoutPicture)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ID token has expired at") {
			return
		}
	}
	assert.Nil(t, err)

	assert.Equal(t, "denis2", id.Username)
	assert.Equal(t, "gY5mfr38daQXLuOxC6mPuz0LS6a2", id.UserId)
	assert.Equal(t, "", id.Picture)
}

func initFirebaseClient() *auth.Client {
	opt := option.WithCredentialsFile("./idp_viewer.json")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v\n", err)
	}

	firebaseClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Could create client: %s", err)
	}
	return firebaseClient
}
