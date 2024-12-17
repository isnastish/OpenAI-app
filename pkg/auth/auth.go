package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/isnastish/openai/pkg/api/models"
)

// NOTE: The subject claim usually identifies one of the parties
// to another (could be user IDs or emails)
// Validation - is the process of checking token's signature.

type Cookie struct {
	Name    string
	Value   string
	Path    string
	Domain  string
	Expires time.Time
	MaxAge  int
}

type AuthManager struct {
	CookieName      string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	JwtSecret       []byte
}

func NewAuthManager(secret []byte) *AuthManager {
	return &AuthManager{
		CookieName:      "__Host-refresh_token",
		AccessTokenTTL:  time.Minute * 15,
		RefreshTokenTTL: time.Hour * 24,
		JwtSecret:       secret,
	}
}

func (a *AuthManager) GetTokenPair(userEmail string) (*models.TokenPair, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&models.Claims{
			// NOTE: There is no need to make the password be a part of claims
			Email: userEmail,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "openai-server",
				// Subject:   "",
				ID:       "1",
				Audience: []string{"openai-frontend"},
			},
		})

	// TODO: Sign a token using a secret key,
	// ideally it should be a private key.
	signedAccessToken, err := token.SignedString(a.JwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %v", err)
	}

	// create a new refresh token with claims
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&models.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				// NOTE: Supposed to be a user ID in a database
				Subject:  userEmail,
				IssuedAt: jwt.NewNumericDate(time.Now()),
				// NOTE: An expiration time for refresh token should be 24 hours,
				// after that a user will be prompted to login again
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.RefreshTokenTTL)),
			},
		},
	)

	// Sign an access token using our secret key.
	signedRefreshToken, err := refreshToken.SignedString(a.JwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign a refresh token: %v", err)
	}

	return &models.TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (a *AuthManager) GetCookie(cookieValue string) *Cookie {
	return &Cookie{
		Name:    a.CookieName,
		Value:   cookieValue,
		Path:    "/",
		Expires: time.Now().Add(a.RefreshTokenTTL),
		MaxAge:  int(a.RefreshTokenTTL.Seconds()),
	}
}
