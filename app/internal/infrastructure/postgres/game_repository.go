package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lock-stock-v2/internal/domain/game/model"
	roomModel "lock-stock-v2/internal/domain/room/model"
	userModel "lock-stock-v2/internal/domain/user/model"
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

func (repo *LockStockGameRepository) FindByUser(user *userModel.User) (*model.LockStockGame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var uid, actionDuration, questionDuration, roomId string
	var createdAt time.Time
	var roomStatus roomModel.RoomStatus

	err := repo.db.QueryRow(ctx, `
	SELECT g.uid, g.action_duration, g.question_duration, g.created_at, r.uid, r.status
	FROM lock_stock_games g
	JOIN players p ON g.id = p.game_id
	JOIN users u ON p.user_id = u.id
	JOIN rooms r ON r.id = g.room_id
	WHERE u.uid = $1
`, user.Uid()).Scan(&uid, &actionDuration, &questionDuration, &createdAt, &roomId, &roomStatus)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	room := roomModel.NewRoom(roomId, roomStatus)
	return model.NewLockStockGame(uid, actionDuration, questionDuration, room), nil
}
