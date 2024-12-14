package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/isnastish/openai/pkg/auth"
	"github.com/isnastish/openai/pkg/db"
	"github.com/isnastish/openai/pkg/db/postgres"
	"github.com/isnastish/openai/pkg/ipresolver"
	"github.com/isnastish/openai/pkg/log"
	"github.com/isnastish/openai/pkg/openai"
)

type UserData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// TODO: This should probably be renamed to server instead of App
type App struct {
	// http server
	fiberApp *fiber.App
	// client for interacting with openai model
	openaiClient *openai.Client
	// ip resovler client for retrieving geolocation data
	ipResolverClient *ipresolver.Client
	// authentication manager
	auth *auth.AuthManager
	// database controller for persisting data
	dbController db.DatabaseController
	// settings
	port int
}

func NewApp(port int /*TODO: pass a secret */) (*App, error) {
	openaiClient, err := openai.NewOpenAIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create an OpenAI client, error: %v", err)
	}

	ipResolverClient, err := ipresolver.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create an ipresolver client, error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// NOTE: For now let's go with db controller,
	// but we should select the contoller based on some env
	// variable DATABASE_CONTROLLER, for example.
	dbContoller, err := postgres.NewPostgresController(ctx)
	if err != nil {
		return nil, err
	}

	app := &App{
		fiberApp: fiber.New(fiber.Config{
			Prefork:      true,
			ServerHeader: "Fiber",
		}),
		openaiClient:     openaiClient,
		ipResolverClient: ipResolverClient,
		auth:             auth.NewAuthManager([]byte("my-dummy-secret")),
		dbController:     dbContoller,
		port:             port,
	}

	// CORS middleware
	app.fiberApp.Use("/", SetupCORSMiddleware)

	// logging middleware
	app.fiberApp.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} latency:${latency} ${status} - ${method} ${path}\n",
	}))

	app.fiberApp.Post("/signup", app.SignupRoute)
	app.fiberApp.Post("/login", app.LoginRoute)
	app.fiberApp.Get("/logout", app.LogoutRoute)
	app.fiberApp.Get("/refresh-token", app.RefreshCookieRoute)
	app.fiberApp.Put("/openai", app.OpenaAIRoute)

	return app, nil
}

func (a *App) Serve() error {
	log.Logger.Info("Listening on port: %v", a.port)

	if err := a.fiberApp.Listen(fmt.Sprintf(":%d", a.port)); err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return nil
}

func (a *App) Shutdown() error {
	// close database connnection first
	defer a.dbController.Close()

	if err := a.fiberApp.Shutdown(); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	return nil
}
