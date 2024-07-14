package main

import (
	"fmt"
	"github.com/fmo/players-api/config"
	"github.com/fmo/players-api/internal/database"
	"github.com/fmo/players-api/internal/services"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type AppConfig struct {
	PlayersService services.PlayersService
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

	db := database.NewDbAdapter()
	playersService := services.NewPlayers(db, logger)

	app := AppConfig{
		PlayersService: playersService,
	}

	srv := &http.Server{
		Addr:    portNumber,
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	log.Fatal(err)
}
