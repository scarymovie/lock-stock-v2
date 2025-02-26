package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var roundID int
	err := repo.db.QueryRow(ctx, "SELECT id FROM rounds WHERE uid = $1", bet.Round().Uid()).Scan(&roundID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("round not found for uid: %s", bet.Round().Uid())
		}
		return fmt.Errorf("failed to find round_id: %w", err)
	}

	query := `
        INSERT INTO bets (
                          amount,
                          player_id,
                          round_id,
                          number)
        VALUES (
					$1,
					(SELECT id FROM players WHERE uid = $2),
					$3, 
					$4
                )
    `

	_, err = repo.db.Exec(ctx, query,
		bet.Amount(),
		bet.Player().Uid(),
		roundID,
		int(bet.Number()),
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
