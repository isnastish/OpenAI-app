package api

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/isnastish/openai/pkg/openai"
)

func (a *App) OpenaAIMessageHandler(ctx *fiber.Ctx) error {
	messsage := ctx.Params("message")

	fmt.Printf("Got a message: %s\n", messsage)

	reqBody := ctx.Request().Body()

	var reqData openai.OpenAIRequest
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
	}

	resp, err := a.OpenaiClient.AskOpenAI(reqData.OpenaiQuestion)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("openai: %s", err.Error()))
	}

	return ctx.JSON(map[string]string{
		"model":  resp.Model,
		"openai": resp.Choices[0].Message.Content,
	}, "application/json")
}

func (a *App) RefreshCookieHanlder(ctx *fiber.Ctx) error {
	// cookieRefreshToken := ctx.Cookies(cookieName)
	// if cookieRefreshToken != "" {
	// 	fmt.Printf("Cookie value: %s\n", cookieRefreshToken)
	// }
	return nil
}

func (a *App) LoginHandler(ctx *fiber.Ctx) error {
	//
	// TODO: Determine where to put user data
	//

	var userData UserData
	if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
	}

	//
	// TODO: Email and password validation, and make sure that such user exists in a database.
	//

	tokenPair, err := a.Auth.GetTokensPair(userData.Email, userData.Password)
	if err != nil {
	}

	cookie := a.Auth.GetCookie(tokenPair.RefreshToken)

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
	// NOTE: In order to delete a cookie we should include
	// the same cookie into a request which contains the same fields
	// with expiry date set to the past, and maxage set to -1
	// cookieRefreshToken := ctx.Cookies(cookieName)
	// if cookieRefreshToken != "" {
	// 	fmt.Printf("Cookie refresh token (logout): %s\n", cookieRefreshToken)
	// }

	// Will remove the cookie on the client side
	// ctx.Cookie(&fiber.Cookie{
	// 	Name:     cookieName,
	// 	Value:    "",
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(-(time.Hour * 2)),
	// 	MaxAge:   -1,
	// 	HTTPOnly: true,
	// 	Secure:   true,
	// 	SameSite: fiber.CookieSameSiteStrictMode,
	// })

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
