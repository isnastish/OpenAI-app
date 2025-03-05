package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/isnastish/aiclient/internal/domain/ipresolver"
	"github.com/isnastish/aiclient/internal/domain/users"
)

type FirestoreUserRepository struct {
	client *firestore.Client
}

func NewFirestoreUserRepository(client *firestore.Client) *FirestoreUserRepository {
	return &FirestoreUserRepository{
		client: client,
	}
}

func (f FirestoreUserRepository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}

func (f FirestoreUserRepository) GetUserByID(ctx context.Context, id int) (*users.User, error) {
	return nil, nil
}

func (f FirestoreUserRepository) InsertUser(ctx context.Context, userData *users.User, geolocation *ipresolver.UserGeolocation) error {
	if _, _, err := f.client.Collection("users").Add(ctx, map[string]interface{}{
		"first_name": userData.FirstName,
		"last_name":  userData.LastName,
		"email":      userData.Email,
		"password":   userData.Password,
		"country":    geolocation.Country,
		"city":       geolocation.City,
	}); err != nil {
		return fmt.Errorf("failed to insert new user, error %v", err)
	}
	return nil
}

func (f FirestoreUserRepository) HasUser(ctx context.Context, email string) (bool, error) {
	// TODO: Implement
	return false, nil
}
