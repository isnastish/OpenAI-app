package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/isnastish/openai/pkg/auth"
	"github.com/isnastish/openai/pkg/log"
	"github.com/isnastish/openai/pkg/openai"
)

//
// NOTE: This should be internal to the package.
//

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type App struct {
	OpenaiClient *openai.Client
	Auth         *auth.AuthManager
	FiberApp     *fiber.App
	Port         int
}

func NewApp(port int) (*App, error) {
	openaiClient, err := openai.NewOpenAIClient()
	if err != nil {
		return nil, fmt.Errorf("")
		// log.Logger.Fatal("Failed to create openai client: %v", err)
	}

	app := &App{
		OpenaiClient: openaiClient,
		Auth:         auth.NewAuthManager([]byte("my-dummy-secret")),
		FiberApp: fiber.New(fiber.Config{
			Prefork:      true,
			ServerHeader: "Fiber",
		}),
		Port: port,
	}

	app.FiberApp.Use("/", SetupCORSMiddleware)

	app.FiberApp.Put("/api/openai", app.OpenaAIMessageRoute)
	app.FiberApp.Post("/api/login", app.LoginRoute)
	app.FiberApp.Post("/api/signup", app.SignupRoute)
	app.FiberApp.Get("/api/logout", app.LogoutRoute)
	app.FiberApp.Get("/api/refresh", app.RefreshCookieRoute)

	return app, nil
}

func (a *App) Serve() error {
	log.Logger.Info("Listening on port: %v", a.Port)

	if err := a.FiberApp.Listen(fmt.Sprintf(":%d", a.Port)); err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return nil
}

func (a *App) Shutdown() error {
	if err := a.FiberApp.Shutdown(); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	return nil
}
