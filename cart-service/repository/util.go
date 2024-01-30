package repository

import (
	"cart-service/logger"
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// Initialize establishes a connection to the Redis server using the provided URI.
// It parses the URI, creates a new Redis client, and stores it in the global 'client' variable.
// If an error occurs during the initialization process, it returns the error.
func Initialize(redisURI string) error {
	// Parse the Redis URI
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return errors.Wrap(err, "Unable to initialize the Redis connection")
	}

	// Create a new Redis client and store it in the global 'client' variable
	client = redis.NewClient(opt)
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
	err := client.Ping(ctx).Err()
	if err != nil {
		logger.Logger.Errorf("Ping - Unable to ping the database - %s", err.Error())
		return errors.Wrap(err, "Unable to ping the database")
	}

	logger.Logger.Debug("Ping - OK")
	return nil
}
