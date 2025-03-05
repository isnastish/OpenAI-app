package server

import "github.com/isnastish/aiclient/internal/app"

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app *app.Application) *HttpServer {
}
