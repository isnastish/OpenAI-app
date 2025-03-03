package app

import (
	"context"
	"net/http"

	"github.com/isnastish/aiclient/internal/ports"
)

// NOTE: The application shouldn't hold any ports, instead it should contain of Commands and Queries
// And we shouldn't keep any handlers here either.
type Application struct {
	userPort ports.UserRepository
}

func (a Application) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse request body here...
	a.LoginController(r.Context())
}

func (a Application) LoginController(ctx context.Context) {
	// check if the user exists
	// Validation should come here...
	a.userPort.GetUserByEmail(ctx, "")
}
