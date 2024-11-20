package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/isnastish/openai/pkg/auth"
	"github.com/isnastish/openai/pkg/log"
	"github.com/isnastish/openai/pkg/openai"
)

// TODO: This should be moved to auth package

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email    string `json:"email"`
	Password string `json:"pwd"`
	jwt.RegisteredClaims
}

type TokensPairs struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	openaiClient, err := openai.NewOpenAIClient()
	if err != nil {
		log.Logger.Fatal("Failed to create openai client: %v", err)
	}

	authManager := auth.NewAuthManager([]byte("my-dummy-secret"))

	// the server which will accept requests from the frontend
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "Fiber",
	})

	// CORS middleware
	app.Use("/", func(ctx *fiber.Ctx) error {
		fmt.Println("Middleware function was triggered")
		// TODO: This has to be moved into a separate function
		// CORS - cross origin request sharing
		// This is the address that our frontend is running on
		ctx.Set("Access-Control-Allow-Origin", "http://localhost:3000")
		ctx.Set("Access-Control-Allow-Credentials", "true")

		if ctx.Method() == "OPTIONS" {
			ctx.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
			ctx.Set("Access-Control-Allow-Headers", "Accept, X-CSRF-Token, Content-Type, Authorization")
			return nil
		}

		return ctx.Next()
	})

	app.Put("/api/openai/:message", func(ctx *fiber.Ctx) error {
		messsage := ctx.Params("message")

		fmt.Printf("Got a message: %s\n", messsage)

		reqBody := ctx.Request().Body()

		var reqData openai.OpenAIRequest
		if err := json.Unmarshal(reqBody, &reqData); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
		}

		resp, err := openaiClient.AskOpenAI(reqData.OpenaiQuestion)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("openai: %s", err.Error()))
		}

		return ctx.JSON(map[string]string{
			"model":  resp.Model,
			"openai": resp.Choices[0].Message.Content,
		}, "application/json")
	})

	app.Get("/api/refresh", func(ctx *fiber.Ctx) error {
		// cookieRefreshToken := ctx.Cookies(cookieName)
		// if cookieRefreshToken != "" {
		// 	fmt.Printf("Cookie value: %s\n", cookieRefreshToken)
		// }

		return nil
	})

	app.Post("/api/login", func(ctx *fiber.Ctx) error {
		var userData UserData
		if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
		}

		// TODO: Email and password validation

		tokenPair, err := authManager.GetTokensPair(userData.Email, userData.Password)
		if err != nil {
		}

		cookie := authManager.GetCookie(tokenPair.RefreshToken)

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
	})

	app.Get("/api/logout", func(ctx *fiber.Ctx) error {
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
	})

	app.Post("/api/signup", func(ctx *fiber.Ctx) error {
		var userData UserData
		if err := json.Unmarshal(ctx.Body(), &userData); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Failed to unmarshal request body")
		}

		// TODO: Make sure that the user doesn't exist
		return fiber.NewError(fiber.StatusNotImplemented, "")
	})

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Logger.Info("Listening on port :3031")

	go func() {
		if err := app.Listen(":3031"); err != nil {
			log.Logger.Fatal("Failed to listen on port :3031 %v", err)
		}
	}()

	<-osSigChan
	if err := app.Shutdown(); err != nil {
		log.Logger.Fatal("Failed to shutdown the server: %v", err)
	}

	os.Exit(0)
}
