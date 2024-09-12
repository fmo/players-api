package ports

import (
	"context"
	"github.com/fmo/players-api/internal/application/core/domain"
)

type DBPort interface {
	FindPlayersByTeamId(ctx context.Context, teamId string) (players []domain.Player, err error)
	FindPlayersById(ctx context.Context, playerId string) (player domain.Player, err error)
}
