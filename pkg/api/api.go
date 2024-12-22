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
	openaiClient, err := openai.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create an OpenAI client, error: %v", err)
	}

	ipResolverClient, err := ipresolver.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create an ipresolver client, error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// NOTE: For now let's go with db controller,
	// but we should select the contoller based on some env
	// variable DATABASE_CONTROLLER, for example.
	dbContoller, err := postgres.NewPostgresController(ctx)
	if err != nil {
		return nil, err
	}

	accessTokenTTL := time.Minute * 15
	app := &App{
		fiberApp: fiber.New(fiber.Config{
			// TODO: Figure out the prefork parameter.
			// Read fiber's documentation.
			Prefork:      false,
			ServerHeader: "Fiber",
		}),
		openaiClient:     openaiClient,
		ipResolverClient: ipResolverClient,
		auth:             auth.NewAuthManager([]byte("my-dummy-secret"), accessTokenTTL),
		dbController:     dbContoller,
		port:             port}

	// CORS middleware
	app.fiberApp.Use("/", SetupCORSMiddleware)

	// logging middleware
	app.fiberApp.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} latency:${latency} ${status} - ${method} ${path}\n",
	}))

	// We need to apply auth middleware only to certain routes.
	// The middleware would be invoked only for routes starting with openai
	app.fiberApp.Use("/protected", func(ctx *fiber.Ctx) error {
		return app.auth.AuthorizationMiddleware(ctx)
	})

	app.fiberApp.Post("/signup", app.SignupRoute)
	app.fiberApp.Post("/login", app.LoginRoute)
	app.fiberApp.Get("/logout", app.LogoutRoute)
	app.fiberApp.Get("/refresh-token", app.RefreshCookieRoute)

	// NOTE: This route should be accessed only if the authentication passes.
	app.fiberApp.Post("/protected/openai", app.OpenaAIRoute)

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

	// TODO: Use ShutdownWithContext instead
	if err := a.fiberApp.Shutdown(); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	return nil
}
