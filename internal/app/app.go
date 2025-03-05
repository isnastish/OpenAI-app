package app

import (
	"github.com/isnastish/aiclient/internal/app/command"
	"github.com/isnastish/aiclient/internal/app/query"
)

type Application struct {
	Commands
	Queries
}

type Commands struct {
	CreateUser command.CreateUserHandler
}

type Queries struct {
	AskAi      query.AskAiHandler
	SignupUser query.SignupUserHandler
}

// NOTE: If we put an Application ctor here, which would invole setting different repositories,
// it will no longer be the application layer.
// We cannot include the initialization of different external ip(s)
// as well as databases here.
// Thus the ctor has to go somewhere, and a reasonable way to put it would be in services/application.go.
