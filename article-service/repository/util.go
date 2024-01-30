package repository

import (
	"article-service/logger"
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Initialize establishes a connection to the MongoDB database using the provided URI.
// It sets a timeout of 5 seconds for the connection attempt and returns an error if the
// connection or ping to the database fails.
func Initialize(mongoDbURI string) error {
	// Set a 5-second timeout for the connection attempt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Configure the MongoDB client with the provided URI
	clientOptions := options.Client().ApplyURI(mongoDbURI)
	var err error

	// Connect to the MongoDB database
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Logger.Errorf("Initialize - Unable to initialize the database connection - %s ", err.Error())
		return errors.Wrap(err, "Unable to initialize the database connection")
	}

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Logger.Errorf("Initialize - Unable to ping the database - %s ", err.Error())
		return errors.Wrap(err, "Unable to ping the database")
	}

	logger.Logger.Info("Initialize - initialization successful")
	return nil
}

// PingDatabase checks the health of the existing MongoDB database connection.
// It sets a timeout of 5 seconds for the ping attempt and returns an error if the
// ping to the database fails.
func PingDatabase() error {
	// Set a 5-second timeout for the ping attempt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to verify the connection
	err := client.Ping(ctx, nil)
	if err != nil {
		logger.Logger.Errorf("Ping - ", err.Error())
		return errors.Wrap(err, "Unable to ping the database")
	}

	logger.Logger.Debug("Ping - OK")
	return nil
}
