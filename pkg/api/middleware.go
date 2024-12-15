package api

import "github.com/gofiber/fiber/v2"

// NOTE: Fiber supports cors config, so use that instead
// We can use cors.New(fiber.config{}) directly in the app class

func SetupCORSMiddleware(ctx *fiber.Ctx) error {
	// http://localhost:3000 is the origin that react front-end is running on.
	ctx.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Set("Access-Control-Allow-Credentials", "true")

	if ctx.Method() == "OPTIONS" {
		ctx.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
		ctx.Set("Access-Control-Allow-Headers", "Accept, X-CSRF-Token, Content-Type, Authorization")
		return nil
	}

	return ctx.Next()
}
