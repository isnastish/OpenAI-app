package services

import (
	"context"

	"github.com/isnastish/aiclient/internal/adapters/ai"
	"github.com/isnastish/aiclient/internal/adapters/database"
	"github.com/isnastish/aiclient/internal/adapters/ipresolver"
	"github.com/isnastish/aiclient/internal/app"
	"github.com/isnastish/aiclient/internal/app/command"
	"github.com/isnastish/aiclient/internal/app/query"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicationEnv struct {
	PostgresURL   string
	OpenAIApiKey  string
	IpflareApiKey string
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
	ipResolverRepo := ipresolver.NewIpflareRespository(env.IpflareApiKey)

	aiRepo := ai.NewOpenAiRepository(env.OpenAIApiKey)

	return &app.Application{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(userRepo, ipResolverRepo),
		},
		Queries: app.Queries{
			AskAi:      query.NewAskAiHandler(aiRepo),
			SignupUser: query.NewSignupUserHandler(userRepo),
		},
	}
}
