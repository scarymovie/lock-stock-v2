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

func (repo *BetRepository) FindByRound(round *model.Round) ([]*model.Bet, error) {
	query := `
        SELECT
            b.player_id,
            b.amount,
            b.number
        FROM bets b
        WHERE b.round_id = (SELECT id FROM rounds WHERE uid = $1)
        ORDER BY b.number ASC
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := repo.db.Query(ctx, query, round.Uid())
	if err != nil {
		return nil, fmt.Errorf("failed to find bets by round: %w", err)
	}
	defer rows.Close()

	var bets []*model.Bet
	for rows.Next() {
		var playerId uint
		var amount int
		var number uint

		err := rows.Scan(&playerId, &amount, &number)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bet row: %w", err)
		}

		player, err := repo.findPlayerById(ctx, playerId)
		if err != nil {
			return nil, fmt.Errorf("failed to find player: %w", err)
		}

		bet := model.NewBet(player, amount, round, number)
		bets = append(bets, bet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over bet rows: %w", err)
	}

	return bets, nil
}

func (repo *BetRepository) findPlayerById(ctx context.Context, playerId uint) (*model.Player, error) {
	var playerUid string
	var userId uint
	var balance int
	var status model.PlayerStatus
	var gameId uint

	err := repo.db.QueryRow(ctx, `
        SELECT p.uid, p.user_id, p.balance, p.status, p.game_id 
        FROM players p 
        WHERE p.id = $1`, playerId,
	).Scan(&playerUid, &userId, &balance, &status, &gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to find player: %w", err)
	}

	user, err := repo.findUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	game, err := repo.findGameById(ctx, gameId)
	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}

	return model.NewPlayer(playerUid, user, balance, status, game), nil
}

func (repo *BetRepository) findUserById(ctx context.Context, userId uint) (*userModel.User, error) {
	var uid string
	var username string

	err := repo.db.QueryRow(ctx, `
        SELECT u.uid, u.name 
        FROM users u 
        WHERE u.id = $1`, userId,
	).Scan(&uid, &username)

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return userModel.NewUser(uid, username), nil
}

func (repo *BetRepository) findGameById(ctx context.Context, gameId uint) (*model.LockStockGame, error) {
	var gameUid string
	var actionDuration string
	var questionDuration string
	var roomId int

	err := repo.db.QueryRow(ctx, `
        SELECT uid, action_duration, question_duration, room_id
        FROM lock_stock_games 
        WHERE id = $1`, gameId,
	).Scan(&gameUid, &actionDuration, &questionDuration, &roomId)

	if err != nil {
		return nil, fmt.Errorf("failed to find game: %w", err)
	}

	room := repo.findRoomById(ctx, uint(roomId))

	return model.NewLockStockGame(gameUid, actionDuration, questionDuration, room), nil
}

func (repo *BetRepository) findRoomById(ctx context.Context, roomUd uint) *roomModel.Room {

	var roomUid string
	var roomStatus roomModel.RoomStatus

	err := repo.db.QueryRow(ctx, `
        SELECT uid, status
        FROM rooms 
        WHERE id = $1`, roomUd,
	).Scan(&roomUid, roomStatus)

	if err != nil {
		return nil
	}

	return roomModel.NewRoom(roomUid, roomStatus)
}
