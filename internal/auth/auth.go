package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	// TODO: We should use user ID instead of email address.
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Auth struct {
	CookieName      string
	DefaultIssuer   string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	JwtSecret       []byte
}

func NewAuth(secret []byte, accessTokenTTL time.Duration) *Auth {
	return &Auth{
		CookieName:      "__host_refresh_token",
		DefaultIssuer:   "openai-server",
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: time.Hour * 48,
		JwtSecret:       secret,
	}
}
