package query

import (
	"context"
	"errors"
	"fmt"

	"github.com/isnastish/aiclient/internal/domain/users"
	"github.com/isnastish/aiclient/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Figure out logging, either use a singleton

type SignupUserHandler struct {
	userRepo ports.UserRepository
}

func NewSignupUserHandler(userRepo ports.UserRepository) *SignupUserHandler {
	return &SignupUserHandler{
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
	user.IsValid(existingUser.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// Unauthorized
			return fmt.Errorf("password does match")
		default:
			// ServerInternalError
			return fmt.Errorf("password validation failed, %v", err)
		}
	}

	// tokens, err := a.auth.GetTokens(userData.Email)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// cookie := a.auth.GetCookie(tokens.RefreshToken)
	return nil
}
