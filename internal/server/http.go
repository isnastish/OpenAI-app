package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/isnastish/aiclient/internal/app"
)

type HttpServer struct {
	app *app.Application
}

func NewHttpServer(app *app.Application) *HttpServer {
	return &HttpServer{
		app: app,
	}
}

func (h HttpServer) RunHttpServerOnAddress(address string) {
	apiRouter := chi.NewRouter()
	_ = apiRouter
}

func (h HttpServer) RunHttpServer() {

}

func (h HttpServer) CreateUser(w *http.ResponseWriter, req *http.Request) {

}

func addCorsMiddleware(router *chi.Mux) {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(corsMiddleware.Handler)
}
