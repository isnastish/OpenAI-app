package firebase

import (
	"context"

	"github.com/isnastish/openai/pkg/api/models"
)

type FirebaseController struct {
}

func NewFirebaseController() (*FirebaseController, error) {
	return nil, nil
}

func (fc *FirebaseController) AddUser(ctx context.Context, userData *models.UserData, geolocationData *models.GeolocationData) error {
	return nil
}

func (fc *FirebaseController) HasUser(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (fc *FirebaseController) Close() error {
	return nil
}
