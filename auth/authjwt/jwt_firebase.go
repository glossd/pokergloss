package authjwt

import (
	"context"
	fauth "firebase.google.com/go/auth"
	"github.com/glossd/pokergloss/auth/authid"
)

func VerifyToken(ctx context.Context, c *fauth.Client, tokenStr string) (*authid.IdentityFull, error) {
	token, err := c.VerifyIDToken(ctx, tokenStr)
	if err != nil {
		return nil, err
	}
	return ExtractIdentityFull(token)
}

func ExtractIdentity(token *fauth.Token) (*authid.Identity, error) {
	// username may not be present if token is anonymous
	var username string
	usernameV, ok := token.Claims["username"]
	if ok {
		username = usernameV.(string)
	}

	var picture string
	pictureClaim, ok := token.Claims["picture"]
	if ok {
		picture = pictureClaim.(string)
	}

	return &authid.Identity{UserId: token.UID, Username: username, Picture: picture}, nil
}

func ExtractIdentityFull(token *fauth.Token) (*authid.IdentityFull, error) {
	id, err := ExtractIdentity(token)
	if err != nil {
		return nil, err
	}

	var provider string
	providerClaim, ok := token.Claims["firebase"].(map[string]interface{})["sign_in_provider"]
	if ok {
		provider = providerClaim.(string)
	}

	var emailVerified bool
	emailVerifiedClaim, ok := token.Claims["email_verified"]
	if ok {
		emailVerified = emailVerifiedClaim.(bool)
	}

	return &authid.IdentityFull{
		Identity: *id,
		Provider: provider,
		EmailVerified: emailVerified,
	}, nil
}
