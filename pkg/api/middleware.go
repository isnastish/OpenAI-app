package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TODO: This has to be moved into auth package.
// Presumably, we can move the whole `authMiddleware` logic there.
// Regarding cors, we could specify those somewhere else using fiber.config.
// and get rid of middleware.go file.
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

func (a *App) AuthMiddleware(ctx *fiber.Ctx) error {
	tokenString, err := getTokenFromHeader(ctx)
	if err != nil {
		return err
	}

	err = a.auth.VerifyJWTToken(*tokenString)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
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
