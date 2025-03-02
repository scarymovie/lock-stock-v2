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
		INSERT INTO players (uid, balance, status, user_id, game_id)
		VALUES (
		        $1,
		        $2,
		        $3, 
		        (SELECT id FROM users WHERE uid = $4), 
		        (SELECT id FROM lock_stock_games WHERE uid = $5)
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query,
		player.Uid(),
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

func (repo *PlayerRepository) FindByUserAndRoom(user *userModel.User, room *roomModel.Room) *model.Player {
	query := `
		SELECT 
		    p.id, 
		    p.uid, 
		    p.balance, 
		    p.status, 
		    p.user_id, 
		    p.game_id,
		    lsg.uid as game_uid
		FROM players p
		JOIN users u ON p.user_id = u.id
		JOIN lock_stock_games lsg ON p.game_id = lsg.id
		WHERE u.uid = $1 AND lsg.room_id = (SELECT id FROM rooms WHERE uid = $2)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var player *model.Player
	var game *model.LockStockGame
	err := repo.db.QueryRow(ctx, query, user.Uid(), room.Uid()).Scan(
		&player.ID,
		&player.Uid,
		&player.Balance,
		&player.Status,
		&player.UserID,
		&player.GameID,
		&game.UID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		log.Printf("Query failed: %v", err)
		return nil
	}
	game.SetRoom(room)
	player.SetGame(game)
	player.SetUser(user)
	return player
}

func (repo *PlayerRepository) FindByUserAndGame(user *userModel.User, game *model.LockStockGame) *model.Player {
	query := `
		SELECT 
		    p.uid, 
		    p.balance, 
		    p.status
			p.user_id,
			p.game_id,
		FROM players p
		JOIN users u ON p.user_id = u.id
		JOIN lock_stock_games lsg ON p.game_id = lsg.id
		WHERE u.uid = $1 AND lsg.uid = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var player *model.Player
	err := repo.db.QueryRow(ctx, query, user.Uid(), game.Uid()).Scan(
		&player.Uid,
		&player.Balance,
		&player.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		log.Printf("Query failed: %v", err)
		return nil
	}
	player.SetGame(game)
	player.SetUser(user)
	return player
}
