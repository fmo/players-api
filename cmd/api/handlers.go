package main

import (
	"errors"
	"github.com/fmo/players-api/internal/helpers"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (app *AppConfig) getSquad(w http.ResponseWriter, r *http.Request) {
	teamId := r.URL.Query().Get("teamId")
	if teamId == "" {
		http.Error(w, "teamId is required", http.StatusBadRequest)
		return
	}

	teamIdInt, _ := strconv.Atoi(teamId)
	players, err := app.PlayersService.FindPlayers(teamIdInt)
	if err != nil {
		log.Debugf("canf find players %v", err)
		helpers.ErrorJSON(w, errors.New("cant find a team"), 404)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.JsonResponse{
		Message: "Players",
		Data:    players,
	})
}

func (app *AppConfig) getPlayer(w http.ResponseWriter, r *http.Request) {
	playerId := chi.URLParam(r, "playerId")
	if playerId == "" {
		http.Error(w, "playerId is required", http.StatusBadRequest)
		return
	}

	player, err := app.PlayersService.FindOnePlayer(playerId)
	if err != nil {
		http.Error(w, "some error happened", http.StatusBadRequest)
		log.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.JsonResponse{
		Message: "Player",
		Data:    player,
	})
}
