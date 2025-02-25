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

type LockStockGameRepository struct {
	db *pgxpool.Pool
}

func NewPostgresGameRepository(db *pgxpool.Pool) *LockStockGameRepository {
	return &LockStockGameRepository{db: db}
}

func (repo *LockStockGameRepository) Save(game *model.LockStockGame) error {
	query := `
		INSERT INTO lock_stock_games (uid, action_duration, question_duration, room_id, created_at)
		VALUES ($1, $2, $3, (SELECT id FROM rooms WHERE uid = $4), $5)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query,
		game.Uid(),
		game.ActionDuration(),
		game.QuestionDuration(),
		game.Room().Uid(),
		game.CreatedAt(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Postgres error: %s, Code: %s, Detail: %s", pgErr.Message, pgErr.Code, pgErr.Detail)
		}
		return fmt.Errorf("failed to save game: %w", err)
	}

	return nil
}
