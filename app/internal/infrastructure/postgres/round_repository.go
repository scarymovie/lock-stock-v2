package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	gameModel "lock-stock-v2/internal/domain/game/model"
	roomModel "lock-stock-v2/internal/domain/room/model"
)

type RoundRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoundRepository(db *pgxpool.Pool) *RoundRepository {
	return &RoundRepository{db: db}
}

func (repo *RoundRepository) FindByGame(game *gameModel.LockStockGame) ([]*gameModel.Round, error) {
	query := `
		SELECT 
			r.uid, r.number, r.buy_in, r.pot,
			g.uid, g.action_duration, g.question_duration,
			rm.uid, rm.status
		FROM rounds r
		JOIN lock_stock_games g ON r.game_id = g.id
		JOIN rooms rm ON g.room_id = rm.id
		WHERE g.uid = $1
		ORDER BY r.number
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := repo.db.Query(ctx, query, game.Uid())
	if err != nil {
		return nil, fmt.Errorf("failed to find rounds: %w", err)
	}
	defer rows.Close()

	var rounds []*gameModel.Round
	for rows.Next() {
		var (
			roundUid         string
			number           uint
			buyIn            uint
			pot              uint
			gameUid          string
			actionDuration   string
			questionDuration string
			roomUid          string
			roomStatus       roomModel.RoomStatus
		)

		if err := rows.Scan(
			&roundUid, &number, &buyIn, &pot,
			&gameUid, &actionDuration, &questionDuration,
			&roomUid, &roomStatus,
		); err != nil {
			return nil, fmt.Errorf("failed to scan round row: %w", err)
		}

		room := roomModel.NewRoom(roomUid, roomStatus)
		game := gameModel.NewLockStockGame(gameUid, actionDuration, questionDuration, room)
		round := gameModel.NewRound(roundUid, &number, buyIn, pot, game)

		rounds = append(rounds, round)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over round rows: %w", err)
	}

	return rounds, nil
}

func (repo *RoundRepository) FindLastByGame(game *gameModel.LockStockGame) (*gameModel.Round, error) {
	if game == nil {
		log.Println("Game is nil")
		return nil, errors.New("game is nil")
	}

	query := `
		SELECT 
			r.uid, r.number, r.buy_in, r.pot,
			g.uid, g.action_duration, g.question_duration,
			rm.uid, rm.status
		FROM rounds r
		JOIN lock_stock_games g ON r.game_id = g.id
		JOIN rooms rm ON g.room_id = rm.id
		WHERE g.uid = $1
		ORDER BY r.number DESC
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := repo.db.QueryRow(ctx, query, game.Uid())

	var (
		roundUid         string
		number           uint
		buyIn            uint
		pot              uint
		gameUid          string
		actionDuration   string
		questionDuration string
		roomUid          string
		roomStatus       roomModel.RoomStatus
	)

	err := row.Scan(
		&roundUid, &number, &buyIn, &pot,
		&gameUid, &actionDuration, &questionDuration,
		&roomUid, &roomStatus,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find last round: %w", err)
	}

	room := roomModel.NewRoom(roomUid, roomStatus)
	game = gameModel.NewLockStockGame(gameUid, actionDuration, questionDuration, room)
	round := gameModel.NewRound(roundUid, &number, buyIn, pot, game)

	return round, nil
}

func (repo *RoundRepository) Save(ctx context.Context, tx pgx.Tx, round *gameModel.Round) error {
	query := `
		INSERT INTO rounds (uid, number, buy_in, pot, game_id)
		VALUES ($1, $2, $3, $4, 
			COALESCE((SELECT id FROM lock_stock_games WHERE uid = $5), -1)
		)
		ON CONFLICT (uid) DO UPDATE 
		SET number = EXCLUDED.number, 
			buy_in = EXCLUDED.buy_in, 
			pot = EXCLUDED.pot
		RETURNING id
	`

	var roundID int
	err := tx.QueryRow(ctx, query,
		round.Uid(),
		round.Number(),
		round.BuyIn(),
		round.Pot(),
		round.Game().Uid(),
	).Scan(&roundID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Postgres error: %s, Code: %s, Detail: %s", pgErr.Message, pgErr.Code, pgErr.Detail)
		}
		return fmt.Errorf("failed to save round: %w", err)
	}

	if roundID == -1 {
		return fmt.Errorf("game with uid %s not found", round.Game().Uid())
	}

	return nil
}
