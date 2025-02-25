package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lock-stock-v2/internal/domain/game/model"
	"log"
	"time"
)

type PlayerRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPlayerRepository(db *pgxpool.Pool) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (repo *PlayerRepository) Save(player *model.Player) error {
	query := `
		INSERT INTO players (balance, status, room_user_id, game_id)
		VALUES ($1, $2, (SELECT id FROM room_users WHERE user_id = (SELECT id FROM users WHERE uid = $3) AND room_id = (SELECT id FROM rooms WHERE uid = $4)), (SELECT id FROM lock_stock_games WHERE uid = $5))
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query,
		player.Balance(),
		player.Status(),
		player.RoomUser().User().Uid(),
		player.Game().Uid(),
		player.Game().Uid(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Postgres error: %s, Code: %s, Detail: %s", pgErr.Message, pgErr.Code, pgErr.Detail)
		}
		return fmt.Errorf("failed to save player: %w", err)
	}

	return nil
}
