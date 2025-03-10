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

type PlayerRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPlayerRepository(db *pgxpool.Pool) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (repo *PlayerRepository) Save(player *model.Player) error {
	query := `
		INSERT INTO players (balance, status, user_id, game_id)
		VALUES (
		        $1,
		        $2, 
		        (SELECT id FROM users WHERE uid = $3), 
		        (SELECT id FROM lock_stock_games WHERE uid = $4)
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query,
		player.Balance(),
		player.Status(),
		player.User().Uid(),
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

func (repo *PlayerRepository) FindByUserAndRoom(user *userModel.User, room *roomModel.Room) (*model.Player, error) {
	query := `
		SELECT 
		    p.balance, 
		    p.status, 
		    lsg.uid as game_uid,
		    lsg.action_duration as game_action_duration,
		    lsg.question_duration as game_question_duration
		FROM players p
		JOIN users u ON p.user_id = u.id
		JOIN lock_stock_games lsg ON p.game_id = lsg.id
		WHERE u.uid = $1 AND lsg.room_id = (SELECT id FROM rooms WHERE uid = $2)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var balance int
	var status model.PlayerStatus
	var gameUid string
	var gameActionDuration string
	var gameQuestionDuration string

	err := repo.db.QueryRow(ctx, query, user.Uid(), room.Uid()).Scan(
		&balance,
		&status,
		&gameUid,
		&gameActionDuration,
		&gameQuestionDuration,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		log.Printf("Query failed: %v", err)
		return nil, err
	}

	game := model.NewLockStockGame(gameUid, gameActionDuration, gameQuestionDuration, room)
	player := model.NewPlayer(user, balance, status, game)

	return player, err
}
