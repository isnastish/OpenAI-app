package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/isnastish/openai/pkg/api/models"
	"github.com/isnastish/openai/pkg/log"
)

// Cloud Firestore creates collections and documents implicitly the first time you add data to the document. You do not need to explicitly create collections or documents.

type FirestoreController struct {
	client *firestore.Client
	// TODO: Either store a collection itself, or its name.
}

type firestoreUserDataWrapper struct {
	FirstName string `firestore:"first_name"`
	LastName  string `firestore:"last_name"`
	Email     string `firestor:"email"`
	Password  string `firestore:"password"`
	Country   string `firestore:"country"`
	City      string `firestore:"city"`
}

func (db *FirestoreController) Close(_ context.Context) error {
	if err := db.client.Close(); err != nil {
		return fmt.Errorf("firestore: failed to close client: %v", err)
	}
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

	// NOTE: firestore documentation doesn't specify if we need to invoke
	// close() method on the client instance.
	databaseId, set := os.LookupEnv("FIRESTORE_DATABASE_ID")
	if set { // the variable is set, but the value might be empty
		client, err = firestore.NewClientWithDatabase(ctx, projectId, databaseId)
	} else {
		client, err = firestore.NewClient(ctx, projectId)
	}

	if err != nil {
		return nil, fmt.Errorf("firesbase: failed to create client: %v", err)
	}

	db := &FirestoreController{
		client: client,
	}

	// NOTE: Just for testing.
	if err := db.AddUser(ctx, &models.UserData{
		FirstName: "Alexey",
		LastName:  "Yevtushenko",
		Email:     "ayevtushenko@gmail.com",
		Password:  "$2b$10$jpYIUC6UjAagw0p0Oh3EzeavyiwqlRsn4KnCYQ4upDTZe4JRLRYZq",
	}, &models.Geolocation{
		Country: "USA",
		City:    "Washington",
	}); err != nil {
		return nil, err
	}

	// Iterate over all the
	iter := db.client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		// TODO: Use DataTo interface
		var wrappedUserData firestoreUserDataWrapper
		if err := doc.DataTo(&wrappedUserData); err != nil {
			return nil, err
		}

		log.Logger.Info("data: %v", wrappedUserData)
	}

	user, err := db.GetUserByEmail(ctx, "ayevtushenko@gmail.com")
	if err != nil {
		return nil, err
	}

	log.Logger.Info("user: %v", user)

	return db, nil
}

func (db *FirestoreController) AddUser(ctx context.Context, userData *models.UserData, geolocation *models.Geolocation) error {
	_, _, err := db.client.Collection("users").Add(ctx, map[string]interface{}{
		"first_name": userData.FirstName,
		"last_name":  userData.LastName,
		"email":      userData.Email,
		"password":   userData.Password,
		"country":    geolocation.Country,
		"city":       geolocation.City,
	})
	if err != nil {
		return fmt.Errorf("firestore: failed to add user, %v", err)
	}

	log.Logger.Info("Successfully added new user")

	return nil
}

func (db *FirestoreController) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	var wrappedUserData firestoreUserDataWrapper

	// TODO: Use WhereEntity instead.
	iter := db.client.Collection("users").Where("email", "==", email).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("firestore: failed to retrieve document, %v", err)
		}
		if err := doc.DataTo(&wrappedUserData); err != nil {
			return nil, fmt.Errorf("firestore: failed to convert document, %v", err)
		}
		break
	}

	// TODO: Retrieve country and city information
	return &models.UserData{
		FirstName: wrappedUserData.FirstName,
		LastName:  wrappedUserData.LastName,
		Email:     wrappedUserData.Email,
		Password:  wrappedUserData.Password,
	}, nil
}

func (db *FirestoreController) GetUserByID(ctx context.Context, id int) (*models.UserData, error) {
	return nil, nil
}
