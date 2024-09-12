package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fmo/players-api/internal/application/core/domain"
	"github.com/fmo/players-api/internal/ports"
	log "github.com/sirupsen/logrus"
	"time"
)

type Application struct {
	cache ports.CachePort
	db    ports.DBPort
}

func NewApplication(cache ports.CachePort, db ports.DBPort) *Application {
	return &Application{
		cache: cache,
		db:    db,
	}
}

func (a Application) Squad(ctx context.Context, teamId string) []domain.Player {
	cacheKey := fmt.Sprintf("squad:%s", teamId)

	squadData, err := a.cache.Get(ctx, cacheKey)
	if err == nil {
		var players []domain.Player
		err = json.Unmarshal([]byte(squadData), &players)
		if err == nil {
			log.Debugf("Found squad for team %s in Redis returning response", teamId)

			return players
		}
	}

	players, err := a.db.FindPlayersByTeamId(ctx, teamId)
	if err != nil {
		log.Debugf("Can't find squad in database for team %s", teamId)

		return players
	}

	jsonSquad, err := json.Marshal(players)
	if err == nil {
		a.cache.Set(ctx, cacheKey, jsonSquad, 240*time.Hour)
	}

	log.Debugf("Found squad for team %s in database returning response", teamId)

	return []domain.Player{}
}

func (a Application) Player(ctx context.Context, playerId string) (domain.Player, error) {
	cacheKey := "player:" + playerId

	playerData, err := a.cache.Get(ctx, cacheKey)

	var player domain.Player

	if err == nil {
		err = json.Unmarshal([]byte(playerData), &player)
		if err == nil {
			return player, nil
		}
	}

	player, err = a.db.FindPlayersById(ctx, playerId)
	if err != nil {
		return domain.Player{}, err
	}

	jsonPlayer, err := json.Marshal(player)
	if err == nil {
		a.cache.Set(ctx, cacheKey, jsonPlayer, 10*time.Minute)
	}

	return player, nil
}
