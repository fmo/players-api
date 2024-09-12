package ports

import (
	"context"
	"github.com/fmo/players-api/internal/application/core/domain"
)

type APIPorts interface {
	Squad(ctx context.Context, teamId int) []domain.Player
	Player(ctx context.Context, playerId string) (domain.Player, error)
}
