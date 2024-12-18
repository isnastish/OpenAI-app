package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isnastish/openai/pkg/api/models"
)

// TODO: Use this prefix instead to retrieve the token value.
const headerPrefix = "Bearer "

// TODO: All these logic should be moved into a separate function,
// probably in auth package.
func (a *App) AuthMiddleware(ctx *fiber.Ctx) error {
	authorizationHeader := ctx.GetRespHeader("authorization")
	if authorizationHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "authorization header missing")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 {
		return fiber.NewError(fiber.StatusUnauthorized, "authorization token is not set")
	}

	if strings.ToLower(headerParts[0]) != "bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "Bearer schema missing")
	}

	claims := models.Claims{}
	_, err := jwt.ParseWithClaims(headerParts[1], claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// NOTE: `alg` key contains a signing method used to sign the JWT token.
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

		// TODO: We have to make sure that an issuer is the same.
		if issuer != a.auth.DefaultIssuer {
			return nil, fmt.Errorf("jwt token invalid, wrong issuer")
		}

		// return secret
		return []byte(a.auth.JwtSecret), nil
	})

	// NOTE: This checks whether a token has expired as well.
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "failed to validate token")
	}

	return ctx.Next()
}

// NOTE: Fiber supports cors config, so use that instead
// We can use cors.New(fiber.config{}) directly in the app class
func SetupCORSMiddleware(ctx *fiber.Ctx) error {
	// http://localhost:3000 is the origin that react front-end is running on.
	ctx.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Set("Access-Control-Allow-Credentials", "true")

	// OPTIONS is the first method when an external origin (our react front-end)
	// makes a request to the back-end server.
	if ctx.Method() == "OPTIONS" {
		ctx.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
		ctx.Set("Access-Control-Allow-Headers", "Accept, X-CSRF-Token, Content-Type, Authorization, x-forwarded-for")
		return nil
	}

	return ctx.Next()
}
