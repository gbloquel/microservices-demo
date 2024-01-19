package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var client *mongo.Client

// Initialize and check the database connection
func Initialize(mongoDbUri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientOptions := options.Client()
	clientOptions.Monitor = otelmongo.NewMonitor()
	clientOptions.ApplyURI(mongoDbUri)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return errors.Wrap(err, "Unable to initialize the database connection")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "Unable to ping the database")
	}

	return nil
}
