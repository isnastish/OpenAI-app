package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/isnastish/openai/pkg/api/models"
	"github.com/isnastish/openai/pkg/openai"
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
	reqBody := ctx.Request().Body()

	var reqData openai.OpenAIRequest
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
	var userData models.UserData
	if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
	}

	tokenPair, err := a.auth.GetTokensPair(userData.Email, userData.Password)
	if err != nil {
	}

	cookie := a.auth.GetCookie(tokenPair.RefreshToken)

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

	return ctx.JSON(tokenPair, "application/json")
}

func (a *App) LogoutRoute(ctx *fiber.Ctx) error {
	// Will remove the cookie on the client side
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
	// Retrieve user's IP address,
	// get geolocation data
	// add user to the database together with its geolocation data
	// set the cookie which contains a session token and a corresponding jwt token?

	var userData models.UserData
	if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
	}

	// TODO: Make sure that the user doesn't exist
	return fiber.NewError(fiber.StatusNotImplemented, "")
}
