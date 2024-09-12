package rest

import (
	"context"
	"errors"
	"github.com/fmo/players-api/internal/adapters/rest/helpers"
	"github.com/fmo/players-api/internal/api"
	"net/http"
	"strconv"
)

func (a Adapter) GetSquad(w http.ResponseWriter, r *http.Request, params api.GetSquadParams) {
	ctx := context.Background()

	players := a.api.Squad(ctx, strconv.Itoa(params.TeamId))

	helpers.WriteJSON(w, http.StatusOK, helpers.JsonResponse{
		Message: "Players",
		Data:    players,
	})
}

func (a Adapter) GetPlayers(w http.ResponseWriter, r *http.Request, playerId string) {
	ctx := context.Background()

	player, err := a.api.Player(ctx, playerId)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("error here"), http.StatusBadRequest)
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.JsonResponse{
		Message: "Player",
		Data:    player,
	})
}
