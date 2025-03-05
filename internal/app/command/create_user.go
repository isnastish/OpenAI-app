package command

import (
	"context"
	"fmt"

	"github.com/isnastish/aiclient/internal/domain/users"
	"github.com/isnastish/aiclient/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserHandler struct {
	userRepo       ports.UserRepository
	ipResolverRepo ports.IpResolverRepository
}

func NewCreateUserHandler(userRepo ports.UserRepository, ipResolverRepo ports.IpResolverRepository) *CreateUserHandler {
	return &CreateUserHandler{
		userRepo:       userRepo,
		ipResolverRepo: ipResolverRepo,
	}
}

func (h CreateUserHandler) Handle(ctx context.Context, user *users.User, ipAddress string) error {
	// TODO: Add email validation using AWS SES service for validating emails.
	existingUser, err := h.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return fmt.Errorf("user [%s] already exists", user.Email)
	}

	location, err := h.ipResolverRepo.GetUserGeolocationData(ctx, ipAddress)
	if err != nil {
		return err
	}
	userPassSha256, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(userPassSha256)
	if err := h.userRepo.InsertUser(ctx, user, location); err != nil {
		return err
	}

	return nil
}
