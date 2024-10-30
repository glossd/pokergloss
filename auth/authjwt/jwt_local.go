package authjwt

import (
	"github.com/dgrijalva/jwt-go"
"github.com/glossd/pokergloss/auth/authid"
)

func VerifyTokenLocalFull(tokenStr string) (*authid.IdentityFull, error) {
	token, err := jwt.Parse(tokenStr, nil)
	if token == nil {
		return nil, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	userId := claims["user_id"].(string)
	var username string
	usernameStr, ok := claims["username"].(string)
	if ok {
		username = usernameStr
	}
	var picture string
	pictureV, ok := claims["picture"]
	if ok {
		picture = pictureV.(string)
	}
	var emailVerified bool
	emailVerifiedClaim, ok := claims["email_verified"]
	if ok {
		emailVerified = emailVerifiedClaim.(bool)
	}

	var provider string
	providerClaim, ok := claims["firebase"].(map[string]interface{})["sign_in_provider"]
	if ok {
		provider = providerClaim.(string)
	}

	return &authid.IdentityFull{
		Identity:      authid.Identity{UserId: userId, Username: username, Picture: picture},
		EmailVerified: emailVerified,
		Provider:      provider,
	}, nil
}
