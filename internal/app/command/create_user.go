package command

import (
	"github.com/isnastish/aiclient/internal/ports"
)

type createUserHandler struct {
	userRepo ports.UserRepository
}

func NewCreateUserHandler(userRepo ports.UserRepository) {

}
