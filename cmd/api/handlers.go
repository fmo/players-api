package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fmo/players-api/internal/api"
	"github.com/fmo/players-api/internal/helpers"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	app AppConfig
}

func NewServer(a AppConfig) Server {
	return Server{
		app: a,
	}
}

func (h Server) GetSquad(w http.ResponseWriter, r *http.Request, params api.GetSquadParams) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("squad:%d", params.TeamId)

	// Try to get squad from Redis
	squadData, err := h.app.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		var players []api.Player
		err = json.Unmarshal([]byte(squadData), &players)
		if err == nil {
			h.app.PlayersService.Logger.Debugf("Found squad for team %d in Redis returning response", params.TeamId)
			helpers.WriteJSON(w, http.StatusOK, helpers.JsonResponse{
				Message: "Players",
				Data:    players,
			})
			return
		}
	}

	// If not found in Redis, get from database
	players, err := h.app.PlayersService.FindPlayersByTeamId(params.TeamId)
	if err != nil {
		h.app.PlayersService.Logger.Debugf("can't find players %v", err)
		helpers.ErrorJSON(w, errors.New("can't find a team"), 404)
		return
	}

	// Store squad in Redis
	jsonSquad, err := json.Marshal(players)
	if err == nil {
		h.app.RedisClient.Set(ctx, redisKey, jsonSquad, 240*time.Hour)
	}

	h.app.PlayersService.Logger.Debugf("Found squad for team id %d in database returning response", params.TeamId)

	helpers.WriteJSON(w, http.StatusOK, helpers.JsonResponse{
		Message: "Players",
		Data:    players,
	})
}

func (h Server) GetPlayers(w http.ResponseWriter, r *http.Request, playerId string) {
	ctx := context.Background()
	redisKey := "player:" + playerId

	// Try to get player from Redis
	playerData, err := h.app.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		var player api.Player
		err = json.Unmarshal([]byte(playerData), &player)
		if err == nil {
			helpers.WriteJSON(w, http.StatusOK, helpers.JsonResponse{
				Message: "Player",
				Data:    player,
			})
			return
		}
	}

	// If not found in Redis, get from database
	player, err := h.app.PlayersService.FindPlayerById(playerId)
	if err != nil {
		http.Error(w, "some error happened", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Store player in Redis
	_, err = json.Marshal(player)
	if err == nil {
		h.app.RedisClient.Set(ctx, redisKey, playerData, 10*time.Minute)
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.JsonResponse{
		Message: "Player",
		Data:    player,
	})
}
