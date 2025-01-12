package mongodb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
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

func (db *MondgodbController) Close(ctx context.Context) error {
	if err := db.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("mongodb: failed to disconnect mongodb client, error: %v", err)
	}
	return nil
}

func NewMongodbController(ctx context.Context) (*MondgodbController, error) {
	mongodbUri, set := os.LookupEnv("MONGODB_URI")
	if !set {
		return nil, fmt.Errorf("MONGODB_URI is not set")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbUri))
	if err != nil {
		return nil, fmt.Errorf("mongodb: failed to create mongodb client, error: %v", err)
	}

	// TODO: Consider startup timeout, that should be specified with client options.
	// if err := client.Ping(ctx, nil); err != nil {
	// 	return nil, fmt.Errorf("mongodb: server is unavailable, error: %v", err)
	// }

	usersCollection := client.Database("users_database").Collection("users")

	return &MondgodbController{
		collection: usersCollection,
		client:     client,
	}, nil
}

func (db *MondgodbController) AddUser(ctx context.Context, userData *models.UserData, geolocation *models.Geolocation) error {
	// NOTE: Omit the geolocation data for now.
	// And, we have to check whether a user with a specified email address already exists.
	// TODO: Hash password together with a salt before adding to a collection.
	// NOTE: We shouldn't keep uer's geolocation data in a database, it doesn't make sense.
	result, err := db.collection.InsertOne(ctx, userData)
	if err != nil {
		return fmt.Errorf("mongodb: failed to add a new user, error: %v", err)
	}

	log.Logger.Info("Adder a new user with ID: %v", result.InsertedID)

	return nil
}

func (db *MondgodbController) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	var result models.UserData
	if err := db.collection.FindOne(ctx, bson.M{"email": email}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			// NOTE: It's not an error if a user is not found,
			// so we just return nil for the user and for the error.
			// Error is returned ONLY when the query failed.
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (db *MondgodbController) GetUserByID(ctx context.Context, id int) (*models.UserData, error) {
	// TODO: We would have to reconsider this function since an id in mongo's collection
	// is represented as a string.
	return nil, nil
}
