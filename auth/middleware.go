package auth

import (
	"context"
	"github.com/gin-gonic/gin"
"github.com/glossd/pokergloss/auth/authconf"
"github.com/glossd/pokergloss/auth/authid"
"github.com/glossd/pokergloss/auth/authjwt"
"github.com/glossd/pokergloss/auth/authunsafe"
"log"
)

const (
	identityKey = "identity"
)

func EmailVerifiedMiddleware(c *gin.Context) {
	token, ok := extractToken(c)
	if !ok {
		noToken401(c)
		return
	}
	identity, err := ParseJwtToken(c.Request.Context(), token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "id token is invalid"})
		return
	}

	if !identity.EmailVerified {
		c.AbortWithStatusJSON(403, gin.H{"error": "email is not verified"})
		return
	}

	c.Set(identityKey, *identity)

	c.Next()
}

func extractToken(c *gin.Context) (token string, ok bool) {
	v := c.Request.Header.Get("Authorization")
	if v == "" || len(v) < len("Bearer ") {
		return "", false
	}

	return v[len("Bearer "):], true
}

func noToken401(c *gin.Context) {
	c.AbortWithStatusJSON(401, gin.H{"error": "id token not provided"})
}

func Middleware(c *gin.Context) {
	token, ok := extractToken(c)
	if !ok {
		noToken401(c)
		return
	}
	identity, err := ParseJwtToken(c.Request.Context(), token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "id token is invalid"})
		return
	}

	c.Set(identityKey, *identity)

	c.Next()
}

func MiddlewareAnonymous(c *gin.Context) {
	token, ok := extractToken(c)
	if !ok {
		noToken401(c)
		return
	}
	identity, err := ParseJwtToken(c.Request.Context(), token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "id token is invalid"})
		return
	}

	if !identity.IsAnonymous() {
		c.AbortWithStatusJSON(403, gin.H{"error": "id token must be anonymous"})
		return
	}

	c.Set(identityKey, *identity)

	c.Next()
}

func WebsocketMiddleware(c *gin.Context) {
	token := c.Request.URL.Query().Get("token")
	if token == "" {
		noToken401(c)
		return
	}

	identity, err := ParseJwtToken(c.Request.Context(), token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "id token is invalid"})
		return
	}

	c.Set(identityKey, *identity)

	c.Next()
}

func ParseJwtToken(ctx context.Context, tokenStr string) (*authid.IdentityFull, error) {
	if authconf.JwtVerificationDisabled() {
		return authjwt.VerifyTokenLocalFull(tokenStr)
	} else {
		return authjwt.VerifyToken(ctx, authunsafe.FirebaseClient, tokenStr)
	}
}

func Id(c *gin.Context) authid.Identity {
	return IdFull(c).Identity
}

func IdFull(c *gin.Context) authid.IdentityFull {
	identity, ok := c.Get(identityKey)
	if !ok {
		// someone forgot to use the Middleware for the HandleFunc
		log.Fatal("ERROR auth.Middleware is not applied to the endpoint")
	}

	return identity.(authid.IdentityFull)
}

func IdSafe(c *gin.Context) (iden *authid.Identity, ok bool) {
	token, ok := extractToken(c)
	if !ok {
		return nil, false
	}
	idenFull, err := ParseJwtToken(c.Request.Context(), token)
	if err != nil {
		return nil, false
	}

	return &idenFull.Identity, true
}
