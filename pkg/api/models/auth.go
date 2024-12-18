package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	// TODO: We should use user ID instead of email address.
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
