package main

import (
	"errors"
	"github.com/fmo/players-api/internal/helpers"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Server struct {
	app AppConfig
}

func NewServer(a AppConfig) Server {
	return Server{
		app: a,
	}
}

func (h Server) GetSquad(w http.ResponseWriter, r *http.Request, params GetSquadParams) {
	teamId := r.URL.Query().Get("teamId")
	if teamId == "" {
		http.Error(w, "teamId is required", http.StatusBadRequest)
		return
	}

	teamIdInt, _ := strconv.Atoi(teamId)
	players, err := h.app.PlayersService.FindPlayersByTeamId(teamIdInt)
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

func (h Server) GetPlayers(w http.ResponseWriter, r *http.Request, playerId string) {
	player, err := h.app.PlayersService.FindPlayerById(playerId)
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
