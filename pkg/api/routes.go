package api

import (
	"bytes"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/isnastish/openai/pkg/api/models"
)

// TODO: There should be a clear separation between routes and
// business logic that is performed in those routes.
// Probably there should be a separte controller responsible for this.
// It should be relatively straightforward to replace the router,
// without doing any modifications for business logic.
//
// This is a controller which contains all the routes that an application
// exposes.

// IMPORTANT:
// Whenever the user wants to access a protected route or resource,
// the user agent should send the JWT,
// typically in the Authorization header using the Bearer schema.
// The content of the header should look like the following:
// Authorization: Bearer <token>

func (a *App) OpenAIRoute(ctx *fiber.Ctx) error {
	result, err := a.openaiController(ctx.Context(), bytes.Clone(ctx.Body()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(map[string]string{
		"openai": result.Choices[0].Message.Content,
	}, "application/json")
}

func (a *App) RefreshTokensRoute(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies(a.auth.CookieName)
	if refreshToken == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "cookie is not set")
	}

	// TODO: This whole token validation process should be moved into a separate function,
	// inside an auth package.

	// NOTE: We should refresh the token a bit before it will be expired,
	// not after, the only problem is how to do that on the cline side.

	// If the signature check passes we could trust the signed data.
	claims := models.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
		// NOTE: This might not work out.
		// tokenClaims := token.Claims.(*models.Claims)
		// Verify the signing method
		// return a single secret we trust
		return []byte(a.auth.JwtSecret), nil
	})

	if err != nil {
		// TODO: Include error message?
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	_ = token

	// This should never fail if we refresh the token a little bit before it expires.
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return fiber.NewError(fiber.StatusUnauthorized, "token has expired")
	}

	// NOTE: claims.Subject will contain a user ID (but in our case for now it's an email address),
	// We should retrive that email address and make a lookup in a database, whether such user
	// exists, and what's more important the token is not expired.

	// This should probably be a use ID.
	userEmail := claims.Subject

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := a.dbController.GetUserByEmail(dbCtx, userEmail)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unknown user")
	}

	tokenPair, err := a.auth.GetTokens(user.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Create a cookie with a new refresh token and set it.
	cookie := a.auth.GetCookie(tokenPair.RefreshToken)

	// Set an actual cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Expires:  cookie.Expires,
		MaxAge:   cookie.MaxAge,
		HTTPOnly: true, // javascript won't have access to this cookie in a web-browser
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return ctx.JSON(tokenPair, "application/json")
}

func (a *App) LoginRoute(ctx *fiber.Ctx) error {
	tokens, cookie, err := a.loginController(ctx.Context(), ctx.Body())
	if err != nil {
		// NOTE: Currently we don't have a way to distinguish between different error codes.
		// Because it can either be an authorization error, or a server internal error.
		// Unauthorized(401) status should be returned when user specifies invalid username (email) or password.
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Expires:  cookie.Expires,
		MaxAge:   cookie.MaxAge,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return ctx.JSON(tokens, "application/json")
}

func (a *App) LogoutRoute(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     a.auth.CookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return ctx.SendStatus(fiber.StatusOK)
}

func (a *App) SignupRoute(ctx *fiber.Ctx) error {
	var ipAddr string
	if (len(ctx.IPs())) > 0 {
		ipAddr = ctx.IPs()[0]
	} else {
		ipAddr = ctx.IP()
	}

	if err := a.signupController(ctx.Context(), ctx.Body(), ipAddr); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
}
