package main

import (
	"context"

	"github.com/isnastish/aiclient/internal/server"
	"github.com/isnastish/aiclient/internal/services"
)

func main() {
	ctx := context.Background()

	app := services.NewApplication(ctx, &services.ApplicationEnv{})

	server.NewHttpServer(app)
}
