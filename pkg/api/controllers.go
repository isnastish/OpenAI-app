package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/isnastish/openai/pkg/api/models"
	"github.com/isnastish/openai/pkg/auth"
	"github.com/isnastish/openai/pkg/log"
	"github.com/isnastish/openai/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Use amazon SES service to authenticate email address.
// NOTE: Controllers should be moved into a separate package.
// The idea behind restructuring is that we want to have a business
// logic fully isolated from underlying database and any web HTTP framework.
// So we can easily switch between those things.
// For example replace fiber with Echo etc.

func unmarshalRequestData[T models.UserData | models.OpenAIRequest](requestBody []byte) (*T, error) {
	var data T
	if err := json.Unmarshal(requestBody, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %v", err)
	}
	return &data, nil
}

func (a *App) openaiController(ctx context.Context, requestBody []byte) (*models.OpenAIResp, error) {
	query, err := unmarshalRequestData[models.OpenAIRequest](requestBody)
	if err != nil {
		return nil, err
	}

	result, err := a.openaiClient.AskOpenAI(ctx, query.OpenaiQuestion)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) loginController(ctx context.Context, requestBody []byte) (*models.Tokens, *auth.Cookie, error) {
	userData, err := unmarshalRequestData[models.UserData](requestBody)
	if err != nil {
		return nil, nil, err
	}

	// Query the user in a database
	existingUser, err := a.dbController.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return nil, nil, err
	}
	if existingUser == nil {
		// Unauthorized
		return nil, nil, fmt.Errorf("user with email: %s doesn't exist", userData.Email)
	}

	// Match password hash
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userData.Password)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// Unauthorized
			return nil, nil, fmt.Errorf("password does match")
		default:
			// ServerInternalError
			return nil, nil, fmt.Errorf("password validation failed, %v", err)
		}
	}

	tokens, err := a.auth.GetTokens(userData.Email)
	if err != nil {
		return nil, nil, err
	}

	cookie := a.auth.GetCookie(tokens.RefreshToken)

	return tokens, cookie, nil
}

func (a *App) signupController(ctx context.Context, requestBody []byte, ipAddr string) error {
	userData, err := unmarshalRequestData[models.UserData](requestBody)
	if err != nil {
		return err
	}

	if err := validator.ValidateUserPassword(userData.Password); err != nil {
		return err
	}

	// TODO: Email validation.
	// a.awsEmailService.SendEmail()

	// Check if the user with given email address already exists.
	existingUser, err := a.dbController.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return err
	}

	// NOTE: Probably internal server error is not the best way of doing this.
	// We should return 409 -> Conflict, or so.
	if existingUser != nil {
		return fmt.Errorf("user %s already exist", userData.Email)
	}

	geolocation, err := a.ipResolverClient.GetGeolocationData(ipAddr)
	if err != nil {
		return fmt.Errorf("faield to get geolocation, %v", err)
	}
	log.Logger.Info("Geolocation: %v", geolocation)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to encrypt password, %v", err)
	}
	userData.Password = string(passwordHash)

	if err := a.dbController.AddUser(ctx, userData, geolocation); err != nil {
		return fmt.Errorf("failed to add user, %v", userData.Email)
	}
	log.Logger.Info("Successfully added a new user")

	return nil
}

func (a *App) refreshTokenController() (*models.Tokens, *auth.Cookie, error) {
	return nil, nil, nil
}
