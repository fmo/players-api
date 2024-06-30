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

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)
}

func main() {
	// if its empty 80 is being used
	portNumber := config.GetApiPort()

	fmt.Println(fmt.Sprintf("Starting app on port %s", portNumber))

	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	db := database.NewDbAdapter()
	playersService := services.NewPlayers(db)

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
