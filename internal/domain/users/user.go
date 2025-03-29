package users

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (u User) IsValid(passSha256 string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passSha256)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return fmt.Errorf("password does match")
		default:
			return fmt.Errorf("password validation failed, error %v", err)
		}
	}
	return nil
}

// TODO: Introduce custom error codes to distinguish between
// bcrypt internal error and passwords not matching.
