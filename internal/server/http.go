package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/isnastish/aiclient/internal/app"
	"github.com/isnastish/aiclient/internal/domain/users"
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

	// Why do we need to set a middleware using apiRouter,
	// and then use a rootRouter for setting the rest?
}

func (h HttpServer) RunHttpServer() {

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

//
// TODO: Figure out how we can get a user from context.
// And maybe move it to handlers.go
//

func (h HttpServer) CreateUser(w http.ResponseWriter, req *http.Request) {
	// NOTE: This should be moved into a middleware since we do the same procedure in each handler.
	var user users.User
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var ip string
	if xffHeader := req.Header.Get("X-Forwarded-For"); xffHeader != "" {
		ip = strings.Split(xffHeader, ",")[0]
	} else {
		ip = strings.Split(req.RemoteAddr, ":")[0]
	}

	err := h.app.Commands.CreateUser.Handle(req.Context(), &user, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	// TODO: Implement authentication by handling cookies.
}

func (h HttpServer) LoginUser(w http.ResponseWriter, req *http.Response) {

}

func (h HttpServer) AskAi(w http.ResponseWriter, req *http.Response) {

}

func (h HttpServer) RefreshToken(w http.ResponseWriter, req *http.Response) {

}
