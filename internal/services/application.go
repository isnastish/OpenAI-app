package services

import (
	"context"

	"github.com/isnastish/aiclient/internal/adapters/database"
	"github.com/isnastish/aiclient/internal/adapters/ipresolver"
	"github.com/isnastish/aiclient/internal/app"
	"github.com/isnastish/aiclient/internal/app/command"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: Move config out of the application construction.
type ApplicationEnv struct {
	PostgresURL string
}

func NewApplication(ctx context.Context, env *ApplicationEnv) *app.Application {
	// TODO: Setup postgres/firestore client here and create a repository
	// NOTE: How do we handle closing the connection?
	postgresConfig, err := pgxpool.ParseConfig(env.PostgresURL)
	if err != nil {
		panic(err)
	}
	postgresConnPool, err := pgxpool.NewWithConfig(ctx, postgresConfig)
	if err != nil {
		panic(err)
	}

	userRepo := database.NewPostgresUserRepository(postgresConnPool)
	ipResolverRepo := ipresolver.NewIpflareRespository()

	return &app.Application{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(userRepo, ipResolverRepo),
		},
	}
}
