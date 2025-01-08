package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"

	"github.com/isnastish/openai/pkg/api/models"
)

// Cloud Firestore stores data in Documents,
// which are stored in Collections.
// Cloud Firestore creates collections and documents implicitly the first time you add data to the document. You do not need to explicitly create collections or documents.

// NOTE: Context should never be a part of a struct,
// but instead passed to each function that needs it.
type FirestoreController struct {
	client *firestore.Client
}

func (db *FirestoreController) Close() error {
	return nil
}

// TODO: Passing context to a constructor might not be the best solution.
// Since it doesn't contain any deadlines nor cancelations.
func NewFirestoreController(ctx context.Context) (*FirestoreController, error) {
	var client *firestore.Client
	var err error

	projectId, set := os.LookupEnv("FIRESTORE_PROJECT_ID")
	if !set || projectId == "" {
		projectId = firestore.DetectProjectID
	}

	databaseId, set := os.LookupEnv("FIRESTORE_DATABASE_ID")
	if set { // the variable is set, but the value might be empty
		client, err = firestore.NewClientWithDatabase(ctx, projectId, databaseId)
	} else {
		client, err = firestore.NewClient(ctx, projectId)
	}

	if err != nil {
		return nil, fmt.Errorf("firesbase: failed to create client: %v", err)
	}

	return &FirestoreController{
		client: client,
	}, nil
}

func (db *FirestoreController) AddUser(ctx context.Context, userData *models.UserData, geolocationData *models.GeolocationData) error {
	return nil
}

func (db *FirestoreController) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	return nil, nil
}

func (db *FirestoreController) GetUserByID(ctx context.Context, id int) (*models.UserData, error) {
	return nil, nil
}
