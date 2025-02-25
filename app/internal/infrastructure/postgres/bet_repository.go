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

type BetRepository struct {
	db *pgxpool.Pool
}

func NewPostgresBetRepository(db *pgxpool.Pool) *BetRepository {
	return &BetRepository{db: db}
}

func (repo *BetRepository) Save(bet *model.Bet) error {

	query := `
        INSERT INTO bets (
                          amount,
                          player_id,
                          round_id,
                          number)
        VALUES (
					$1,
					(SELECT id FROM players WHERE uid = $2),
					(SELECT id FROM rounds WHERE uid = $3), 
					$4
                )
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query,
		bet.Amount(),
		bet.Player().Uid(),
		bet.Round().Uid(),
		bet.Number(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Postgres error: %s, Code: %s, Detail: %s", pgErr.Message, pgErr.Code, pgErr.Detail)
		}

		return fmt.Errorf("failed to save bet: %w", err)
	}

	return nil
}
