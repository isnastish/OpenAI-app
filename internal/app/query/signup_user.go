package query

import (
	"context"
	"fmt"

	"github.com/isnastish/aiclient/internal/domain/users"
	"github.com/isnastish/aiclient/internal/ports"
)

// TODO: Figure out logging, either use a singleton

type SignupUserHandler struct {
	userRepo ports.UserRepository
}

func NewSignupUserHandler(userRepo ports.UserRepository) SignupUserHandler {
	return SignupUserHandler{
		userRepo: userRepo,
	}
}

// NOTE: This has to be rewritten, instead of accepting a user model,
// what we have to do is to pass a command handler.
// This should be a callback, otherwise it doesn't make any sense to have the same signature.
func (h SignupUserHandler) Handle(ctx context.Context, user *users.User) error {
	existingUser, err := h.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return fmt.Errorf("user [%s] doesn't exist", user.Email)
	}
	if err := user.IsValid(existingUser.Password); err != nil {
		return err
	}

	// NOTE: We shouldn't include the logic for handling cookies here,
	// since it's not a business logic overall.
	// Probaby all of that logic should be inside the http handler.

	// tokens, err := a.auth.GetTokens(userData.Email)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// cookie := a.auth.GetCookie(tokens.RefreshToken)
	return nil
}
