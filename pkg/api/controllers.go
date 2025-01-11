package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/isnastish/openai/pkg/api/models"
	"github.com/isnastish/openai/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Use amazon SES service to authenticate email address.
// NOTE: Controllers should be moved into a separate package.
// The idea behind restructuring is that we want to have a business
// logic fully isolated from underlying database and any web HTTP framework.
// So we can easily switch between those things.
// For example replace fiber with Echo etc.

func unmarshalRequestData[T any](requestBody []byte) (*T, error) {
	var data T
	if err := json.Unmarshal(requestBody, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %v", err)
	}
	return &data, nil
}

func (a *App) openaiRouteImpl(ctx context.Context, requestBody []byte) (*models.OpenAIResp, error) {
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

func (a *App) LoginImpl(ctx context.Context, requestBody []byte) (*models.TokenPair, *auth.Cookie, error) {
	// TODO: We should validate users data,
	// an email address and user's password.
	// In order to do that, we would have to retrieve a user from the database
	// The only problem is that our database contains other data
	// than UserData, its geolocation as well.
	// If passwords don't match we return BadRequest, otherwise
	// we proceed and update access and refresh tokens.
	var userData models.UserData
	if err := json.Unmarshal(requestBody, &userData); err != nil {
		// server internal error or bad request
		return nil, nil, fmt.Errorf("failed to unmarshal request body: %v", err)
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

func (a *App) signupRouteImpl() {

}
