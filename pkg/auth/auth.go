package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isnastish/openai/pkg/api/models"
)

// NOTE: Why the logic for validating a refresh token and jwt token should
// be different?

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
	DefaultIssuer   string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	JwtSecret       []byte
}

func NewAuthManager(secret []byte, accessTokenTTL time.Duration) *AuthManager {
	return &AuthManager{
		CookieName:      "__refresh_token",
		DefaultIssuer:   "openai-server",
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: time.Hour * 48,
		JwtSecret:       secret,
	}
}

func (a *AuthManager) GetTokenPair(userEmail string) (*models.TokenPair, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&models.Claims{
			Email: userEmail,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.AccessTokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    a.DefaultIssuer,
			},
		})

	signedAccessToken, err := token.SignedString(a.JwtSecret)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to sign access token: %v", err)
	}

	// create a new refresh token with claims
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&models.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				// NOTE: Supposed to be a user ID in a database
				Subject:   userEmail,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.RefreshTokenTTL)),
			},
		},
	)

	signedRefreshToken, err := refreshToken.SignedString(a.JwtSecret)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to sign a refresh token: %v", err)
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

func (a *AuthManager) ValidateJwtToken(tokenString string) error {
	claims := models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		tokenClaims := token.Claims.(*models.Claims)
		if tokenClaims.ExpiresAt == nil || time.Now().After(tokenClaims.ExpiresAt.Time) {
			return nil, fmt.Errorf("jwt token is expired")
		}

		issuer, err := tokenClaims.GetIssuer()
		if err != nil {
			return nil, fmt.Errorf("jwt token invlaid, excepted an issuer")
		}

		if issuer != a.DefaultIssuer {
			return nil, fmt.Errorf("jwt token invalid, wrong issuer")
		}

		return []byte(a.JwtSecret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("jwt token is invalid")
	}

	return nil
}

const headerPrefix = "Bearer "

func getTokenFromHeader(ctx *fiber.Ctx) (*string, error) {
	authHeaders, ok := ctx.GetReqHeaders()["Authorization"]
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "authorization header missing")
	}

	for _, header := range authHeaders {
		auth := strings.Clone(header)
		if strings.HasPrefix(auth, headerPrefix) {
			token := strings.TrimSpace(auth[len(headerPrefix):])
			return &token, nil
		}
	}

	return nil, fiber.NewError(fiber.StatusUnauthorized, "authorization header invalid")
}

func (a *AuthManager) AuthorizationMiddleware(ctx *fiber.Ctx) error {
	tokenString, err := getTokenFromHeader(ctx)
	if err != nil {
		return err
	}

	err = a.ValidateJwtToken(*tokenString)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return ctx.Next()
}
