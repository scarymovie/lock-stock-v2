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

type RoundRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoundRepository(db *pgxpool.Pool) *RoundRepository {
	return &RoundRepository{db: db}
}

func (repo *RoundRepository) FindByGame(game *model.LockStockGame) ([]*model.Round, error) {
	query := `
        SELECT
            r.id,
            r.round_number,
            r.question,
            r.price,
            r.game_id
        FROM rounds r
        WHERE r.game_id = (SELECT id FROM lock_stock_games WHERE uid = $1)
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := repo.db.Query(ctx, query, game.Uid())
	if err != nil {
		return nil, fmt.Errorf("failed to find rounds by game: %w", err)
	}
	defer rows.Close()

	var rounds []*model.Round
	for rows.Next() {
		var round model.Round
		var gameID uint

		err := rows.Scan(
			&round.ID,
			&round.RoundNumber,
			&round.QuestionText,
			&round.Price,
			&gameID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan round row: %w", err)
		}
		round.SetGame(game)

		rounds = append(rounds, &round)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over round rows: %w", err)
	}

	return rounds, nil
}

func (repo *RoundRepository) Save(round *model.Round) error {
	query := `
        INSERT INTO rounds (round_number, question, price, game_id)
        VALUES ($1, $2, $3, (SELECT id FROM lock_stock_games WHERE uid = $4))
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// @todo questions
	_, err := repo.db.Exec(ctx, query,
		round.Number(),
		round.BuyIn(),
		round.Game().Uid(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Postgres error: %s, Code: %s, Detail: %s", pgErr.Message, pgErr.Code, pgErr.Detail)
		}
		return fmt.Errorf("failed to save round: %w", err)
	}

	return nil
}
