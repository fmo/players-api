package main

import (
	"context"
	"github.com/fmo/players-api/config"
	"github.com/fmo/players-api/internal/adapters/cache/redis"
	"github.com/fmo/players-api/internal/adapters/db/dynamodb"
	"github.com/fmo/players-api/internal/adapters/rest"
	"github.com/fmo/players-api/internal/application/core/api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

var logger = log.New()

func init() {
	logger.Out = os.Stdout

	logger.Level = log.InfoLevel
}

func main() {
	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	ctx := context.Background()

	cacheAdapter, err := redis.NewAdapter(config.GetRedisAddr(), config.GetRedisPassword())
	if err != nil {
		log.Fatalf("Failed to connect to redis. Error: %v", err)
	}

	dbAdapter, err := dynamodb.NewAdapter(config.GetDynamoDbTableName())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	application := api.NewApplication(cacheAdapter, dbAdapter)
	restAdapter := rest.NewAdapter(application, config.GetApplicationPort())
	restAdapter.Run(ctx)
}
