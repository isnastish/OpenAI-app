package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"github.com/isnastish/openai/pkg/api/models"
)

type FirebaseController struct {
	// firebase app
	app *firebase.App
}

func NewFirebaseController(ctx context.Context) (*FirebaseController, error) {
	// TODO: NewApp should accept the firebase config
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("firebase: failed to create firebase app, error: %v", err)
	}

	return &FirebaseController{
		app: app,
	}, nil
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
