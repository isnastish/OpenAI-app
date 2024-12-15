package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email    string `json:"email"`
	Password string `json:"pwd"`
	jwt.RegisteredClaims
}

type TokensPairs struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
