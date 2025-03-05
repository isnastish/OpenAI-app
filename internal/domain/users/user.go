package users

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
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
