package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email    string `json:"email"`
	Password string `json:"pwd"`
	jwt.RegisteredClaims
}

type TokensPairs struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

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
	JWTSecret       []byte // or a string, but ideally it should be a private key
}

func NewAuthManager(jwtSecret []byte) *AuthManager {
	return &AuthManager{
		CookieName:      "__Host-refresh_token",
		AccessTokenTTL:  time.Minute * 15,
		RefreshTokenTTL: time.Hour * 24,
		JWTSecret:       jwtSecret,
	}
}

func (a *AuthManager) GetTokensPair(userEmailAddress, userPassword string) (*TokensPairs, error) {
	// create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&Claims{
			Email:    userEmailAddress,
			Password: userPassword,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "test",
				// Subject:   "",
				ID:       "1",
				Audience: []string{"openai-frontend"},
			},
		})

	// Sign token using a secret key, it should be private key
	signedAccessToken, err := token.SignedString(a.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %v", err)
	}

	// create a new refresh token with claims
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				// NOTE: Supposed to be a user ID in a database
				// Subject:  userData.Email,
				IssuedAt: jwt.NewNumericDate(time.Now()),
				// NOTE: An expiration time for refresh token should be 24 hours,
				// after that a user will be prompted to login again
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.RefreshTokenTTL)),
			},
		},
	)

	signedRefreshToken, err := refreshToken.SignedString(a.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign a refresh token: %v", err)
	}

	return &TokensPairs{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (a *AuthManager) GetCookie(cookieValue string) *Cookie {
	// TODO: Include domain name
	return &Cookie{
		Name:    a.CookieName,
		Value:   cookieValue,
		Path:    "/",
		Expires: time.Now().Add(a.RefreshTokenTTL),
		MaxAge:  int(a.RefreshTokenTTL.Seconds()),
	}
}
