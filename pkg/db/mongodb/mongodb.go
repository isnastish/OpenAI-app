package mongodb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/isnastish/openai/pkg/api/models"
	"github.com/isnastish/openai/pkg/log"
	_ "github.com/isnastish/openai/pkg/log"
)

type MondgodbController struct {
	// mongodb client
	client *mongo.Client
	// TODO: Do we need to store a collection?
	// Most likely we only need a client.
	collection *mongo.Collection
}

func NewMongodbController(ctx context.Context) (*MondgodbController, error) {
	mongodbUri, set := os.LookupEnv("MONGODB_URI")
	if !set || mongodbUri == "" {
		return nil, fmt.Errorf("MONGODB_URI is not set")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbUri))
	if err != nil {
		return nil, fmt.Errorf("mongodb: failed to create mongodb client, error: %v", err)
	}

	// NOTE: `newdatabase` is a testing database which is not meant to use for production.
	// as well as `users` collection inside that database.
	usersCollection := client.Database("newdatabase").Collection("users")

	return &MondgodbController{
		collection: usersCollection,
		client:     client,
	}, nil
}

func (db *MondgodbController) AddUser(ctx context.Context, userData *models.UserData, geolocationData *models.GeolocationData) error {
	// NOTE: Omit the geolocation data for now.
	// And, we have to check whether a user with a specified email address already exists.
	// TODO: Hash password together with a salt before adding to a collection.
	result, err := db.collection.InsertOne(ctx, userData)
	if err != nil {
		return fmt.Errorf("mongodb: failed to add new user, error: %v", err)
	}

	log.Logger.Info("Adder a new user with ID: %v", result.InsertedID)

	return nil
}

func (db *MondgodbController) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	// db.collection.FindOne(ctx, )
	return nil, nil
}

func (db *MondgodbController) GetUserByID(ctx context.Context, id int) (*models.UserData, error) {
	return nil, nil
}

func (db *MondgodbController) Close() error {
	// TODO: Switch to normal context, but that would probably
	// require refactoring in all db implementations.
	if err := db.client.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("mongodb: failed to disconnect mongodb client, error: %v", err)
	}
	return nil
}
