package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/isnastish/openai/pkg/openai"
)

// NOTE: The business logic shouldn't be in routes, it should be moved to controllers
// Probably in controllers

func (a *App) OpenaAIMessageRoute(ctx *fiber.Ctx) error {
	reqBody := ctx.Request().Body()

	var reqData openai.OpenAIRequest
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
	}

	// make a requests to OpenAI server

	reqCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel() // NOTE: Not sure whether this is the best place where to put this

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
	var userData UserData
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
	var userData UserData
	if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
	}

	// TODO: Make sure that the user doesn't exist
	return fiber.NewError(fiber.StatusNotImplemented, "")
}
