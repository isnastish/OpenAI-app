package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
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

func (a *App) OpenaAIRoute(ctx *fiber.Ctx) error {
	// NOTE: This route should be protected
	// We should validate the token received from the client
	reqBody := ctx.Request().Body()

	var reqData models.OpenAIRequest
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
	}

	reqCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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

func (a *App) RefreshCookieRoute(ctx *fiber.Ctx) error {
	// cookieRefreshToken := ctx.Cookies(cookieName)
	// if cookieRefreshToken != "" {
	// 	fmt.Printf("Cookie value: %s\n", cookieRefreshToken)
	// }
	return nil
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

	tokenPair, err := a.auth.GetTokenPair(userData.Email, userData.Password)
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
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to unmarshal request body, error: %v", err))
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exists, err := a.dbController.HasUser(dbCtx, userData.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if exists {
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

	if err := a.dbController.AddUser(dbCtx, &userData, geolocationData); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to add user, error: %v", err))
	}

	log.Logger.Info("Successfully added user to the database")

	// TODO: Make sure that the user doesn't exist
	return fiber.NewError(fiber.StatusNotImplemented, "")
}
