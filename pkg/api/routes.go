package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/isnastish/openai/pkg/api/models"
	"github.com/isnastish/openai/pkg/log"
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

func (a *App) OpenaAIRoute(ctx *fiber.Ctx) error {
	reqBody := ctx.Request().Body()

	var reqData models.OpenAIRequest
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
	}

	reqCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := a.openaiClient.AskOpenAI(reqCtx, reqData.OpenaiQuestion)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("openai: %s", err.Error()))
	}

	return ctx.JSON(map[string]string{
		"model":  resp.Model,
		"openai": resp.Choices[0].Message.Content,
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

	tokenPair, err := a.auth.GetTokenPair(user.Email)
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
	// TODO: We should validate users data,
	// an email address and user's password.
	// In order to do that, we would have to retrieve a user from the database
	// The only problem is that our database contains other data
	// than UserData, its geolocation as well.
	// If passwords don't match we return BadRequest, otherwise
	// we proceed and update access and refresh tokens.

	var userData models.UserData
	if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to unmarshal request body, error: %v", err))
	}

	// TODO: Perform data validation before making queries to the database.
	// password and email address (on the submitted data).

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := a.dbController.GetUserByEmail(dbCtx, userData.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("user with email: %s doesn't exist", userData.Email))
	}

	// Compare user's password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return fiber.NewError(fiber.StatusBadRequest, "invalid password")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to validate user's password, error: %v", err))
		}
	}

	tokenPair, err := a.auth.GetTokenPair(userData.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	// The cookie holds a refresh token with exactly the same
	// TTL as the refresh token itself.
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

func (a *App) LogoutRoute(ctx *fiber.Ctx) error {
	// Delete the cookie
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
	var userData models.UserData
	if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
		log.Logger.Info("failed to unmarshal request body: error %v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to unmarshal request body, error: %v", err))
	}

	log.Logger.Info("user data: %v", userData)

	dbCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if the user with given email address already exists.
	user, err := a.dbController.GetUserByEmail(dbCtx, userData.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// NOTE: Probably internal server error is not the best way of doing this.
	// We should return 409 -> Conflict, or so.
	if user != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("user with email: %s already exists", userData.Email))
	}

	// Get IP addresses in X-Forwarded-For header
	var ipAddr string
	if (len(ctx.IPs())) > 0 {
		ipAddr = ctx.IPs()[0]
	} else {
		ipAddr = ctx.IP()
	}

	geolocationData, err := a.ipResolverClient.GetGeolocationData(ipAddr)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to retrieve geolocation data, error: %v", err))
	}

	log.Logger.Info("Retrieved geolocation data: %v", geolocationData)

	// TODO: Implement data validation here as well.
	// Move the logic for validating data into a separate function together
	// with encrypting the password.

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("postgres: failed to encrypt password, error: %v", err)
	}

	log.Logger.Info("Encrypted password: %s", hashedPassword)

	userData.Password = string(hashedPassword)

	if err := a.dbController.AddUser(dbCtx, &userData, geolocationData); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to add user, error: %v", err))
	}

	log.Logger.Info("Successfully added user to the database")

	return ctx.SendStatus(fiber.StatusOK)
}
