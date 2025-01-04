package firebase

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"

	"github.com/isnastish/openai/pkg/api/models"
)

type FirestoreController struct {
	// firestore database client
	dbClient *db.Client
}

func NewFirestoreController(ctx context.Context) (*FirestoreController, error) {
	// TODO: Read about refresh-token credentials file
	opt := option.WithCredentialsFile("path/to/refreshToken.json")
	config := &firebase.Config{ProjectID: "my-project-id"}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, fmt.Errorf("firebase: failed to create firebase app, error: %v", err)
	}

	firestoreDbUrl, set := os.LookupEnv("FIRESTORE_URL")
	if !set || firestoreDbUrl == "" {
		return nil, fmt.Errorf("firebase: FIRESTORE_URL is not set")
	}

	client, err := app.DatabaseWithURL(ctx, firestoreDbUrl)
	if err != nil {
		return nil, fmt.Errorf("firebase: failed to initialize database client, error: %v", err)
	}

	// usersRef := client.NewRef("/users")

	return &FirestoreController{
		dbClient: client,
	}, nil
}

func (db *FirestoreController) AddUser(ctx context.Context, userData *models.UserData, geolocationData *models.GeolocationData) error {
	// NOTE: This is a rough approximation of how this supposed to
	usersRef := db.dbClient.NewRef("/users")
	newUserRef, err := usersRef.Push(ctx, nil)
	if err != nil {
		return fmt.Errorf("firestore: failed to create user ref, error: %v", err)
	}

	_ = newUserRef

	return nil
}

func (db *FirestoreController) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	return nil, nil
}

func (db *FirestoreController) GetUserByID(ctx context.Context, id int) (*models.UserData, error) {
	return nil, nil
}

func (db *FirestoreController) Close() error {
	return nil
}
