package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/fmo/football-proto/golang/player"
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

	msgNumber := 0

	for {
		message, err := k.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Errorf("Error reading message: %v\n", err)
			continue
		}

		msgNumber++

		log.Debugf("received the %d. message payload", msgNumber)

		if len(message.Value) == 0 {
			log.Debugf("received an empty message")
			continue
		}

		var players []*pb.Player
		err = json.Unmarshal(message.Value, &players)
		if err != nil {
			log.Errorf("Error unmarshalling message: %v", err)
			continue
		}

		debugPrint := ""
		for _, player := range players {
			debugPrint = fmt.Sprintf("%s, %s", debugPrint, player.Name)
		}

		log.WithFields(log.Fields{
			"partialMessage": fmt.Sprintf("%s...", debugPrint[2:100]),
		}).Debugf("unmarshalled the %d. message payload", msgNumber)

		for _, player := range players {
			if player.Photo != "" {
				playerPhotoName := fmt.Sprintf("%s.png", player.RapidApiId)
				err = s3Service.Save(playerPhotoName, player.Photo)
				if err != nil {
					log.Error(err)
				}
			}

			p := models.Player{
				Team:        player.Team,
				TeamId:      player.TeamId,
				Name:        player.Name,
				Firstname:   player.Firstname,
				Lastname:    player.Lastname,
				Age:         player.Age,
				Nationality: player.Nationality,
				Photo:       player.Photo,
				RapidApiID:  player.RapidApiId,
				Appearances: player.Appearances,
				Position:    player.Position,
			}

			_, err := playersService.CreateOrUpdate(p)
			if err != nil {
				log.Fatalf("Got error calling PutItem: %s", err)
			}

			log.WithFields(log.Fields{
				"playerId":   p.RapidApiID,
				"playerName": p.Name,
				"teamName":   p.Team,
			}).Debug("inserted or updated the player to the database")
		}
	}
}
