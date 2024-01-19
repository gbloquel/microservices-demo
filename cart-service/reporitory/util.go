package reporitory

import (
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func Initialize(redisUri string) error {
	opt, err := redis.ParseURL(redisUri)
	if err != nil {
		return errors.Wrap(err, "Unable to initialize the redis connection")
	}

	client = redis.NewClient(opt)
	// check redis if is ok
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return errors.Wrap(err, "Unable to ping REDIS")
	}

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(client); err != nil {

		return errors.Wrap(err, "Unable to instrumentTracing REDIS")
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(client); err != nil {
		return errors.Wrap(err, "Unable to instrumentMetrics REDIS")
	}
	return nil
}
