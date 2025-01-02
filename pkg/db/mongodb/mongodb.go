package mongodb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/isnastish/openai/pkg/api/models"
	_ "github.com/isnastish/openai/pkg/log"
)

type MondgodbController struct {
	client *mongo.Client
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

	return &MondgodbController{
		client: client,
	}, nil
}

func (db *MondgodbController) AddUser(ctx context.Context, userData *models.UserData, geolocationData *models.GeolocationData) error {
	return nil
}

func (db *MondgodbController) GetUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
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
