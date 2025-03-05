package app

type Application struct {
	// userPort ports.UserRepository
	Commands
	Queries
}

type Commands struct {
}

type Queries struct {
}

// func (a Application) LoginController(ctx context.Context) {
// 	// check if the user exists
// 	// Validation should come here...
// 	a.userPort.GetUserByEmail(ctx, "")
// }
