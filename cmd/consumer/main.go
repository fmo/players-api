package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fmo/players-api/internal/database"
	"github.com/fmo/players-api/internal/kafka"
	"github.com/fmo/players-api/internal/models"
	"github.com/fmo/players-api/internal/s3"
	"github.com/fmo/players-api/internal/services"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)
}

func main() {
	fmt.Println(fmt.Sprintf("Player consumer is up and running"))

	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	k := kafka.NewKafka()
	defer k.Reader.Close()

	playersService := services.NewPlayers(database.NewDbAdapter())

	s3Service, err := s3.NewS3Service()
	if err != nil {
		log.Fatalf("cant connect to s3 %v", err)
	}

	for {
		message, err := k.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Errorf("Error reading message: %v\n", err)
			continue
		}

		log.Debug("Received the payload")

		if len(message.Value) == 0 {
			log.Debugf("Received an empty message")
			continue
		}

		var players []models.Player
		err = json.Unmarshal(message.Value, &players)
		if err != nil {
			log.Errorf("Error unmarshalling message: %v", err)
			continue
		}

		log.WithFields(log.Fields{
			"unmarshalled-payload": players,
		}).Debug("Unmarshalled the payload")

		for _, player := range players {
			if player.Photo != "" {
				playerPhotoName := fmt.Sprintf("%s.png", player.RapidApiID)
				err = s3Service.Save(playerPhotoName, player.Photo)
				if err != nil {
					log.Error(err)
				}
			}

			itemOutput, err := playersService.CreateOrUpdate(player)
			if err != nil {
				log.Fatalf("Got error calling PutItem: %s", err)
			}

			log.WithFields(log.Fields{
				"insertedItemOutput": itemOutput.String(),
			}).Debug("Inserted or updated the player to the database")
		}
	}
}
