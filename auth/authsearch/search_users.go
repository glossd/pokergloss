package authsearch

import (
	"context"
	"firebase.google.com/go/auth"
"github.com/glossd/pokergloss/auth/authconf"
"github.com/glossd/pokergloss/auth/authid"
"github.com/glossd/pokergloss/auth/authunsafe"
)

func GetIdentities(ctx context.Context, userIDs []string) ([]authid.Identity, error) {
	if authconf.JwtVerificationDisabled() {
		var identities []authid.Identity
		for _, uid := range userIDs {
			identities = append(identities, authid.Identity{UserId: uid, Username: "mock username"})
		}
		return identities, nil
	}

	var identifiers []auth.UserIdentifier
	for _, uid := range userIDs {
		identifiers = append(identifiers, auth.UIDIdentifier{UID: uid})
	}
	records, err := authunsafe.FirebaseClient.GetUsers(ctx, identifiers)
	if err != nil {
		return nil, err
	}

	var idens []authid.Identity
	for _, record := range records.Users {
		iden, err := authid.FromRecord(record)
		if err != nil {
			return nil, err
		}
		idens = append(idens, *iden)
	}

	return idens, nil
}
