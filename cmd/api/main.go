package main

import (
	"context"
	"fmt"
	"github.com/fmo/players-api/config"
	"github.com/fmo/players-api/internal/database"
	"github.com/fmo/players-api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type AppConfig struct {
	PlayersService services.PlayersService
	RedisClient    *redis.Client
}

var logger = log.New()

func init() {
	logger.Out = os.Stdout

	logger.Level = log.DebugLevel
}

func main() {
	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// if its empty 80 is being used
	portNumber := config.GetApiPort()
	fmt.Println(fmt.Sprintf("Starting app on port %s", portNumber))

	// initiate database
	db := database.NewDbAdapter()

	// create player service
	playersService := services.NewPlayers(db, logger)

	// connect to Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	conn, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	} else {
		log.Debugf("Connected to redis %v", conn)
	}

	// define new server and assign app config
	server := NewServer(AppConfig{
		PlayersService: playersService,
		RedisClient:    redisClient,
	})

	r := chi.NewMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(corsHandler.Handler)

	h := HandlerFromMux(server, r)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: h,
	}

	log.Fatal(srv.ListenAndServe())
}
