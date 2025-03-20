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

var FailedOnInsertingData = errors.New("failed on inserting data")

func (repo *BetRepository) Save(ctx context.Context, tx pgx.Tx, bet *model.Bet) error {
	query := `
	INSERT INTO bets (
		amount,
		player_id,
		round_id
	)
	SELECT 
		$1,
		p.id,
		r.id
	FROM players p
	JOIN users u ON p.user_id = u.id
	JOIN rounds r ON r.uid = $2
	WHERE u.uid = $3
	`

	result, err := tx.Exec(ctx, query,
		bet.Amount(),
		bet.Round().Uid(),
		bet.Player().User().Uid(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Postgres error: %s, Code: %s, Detail: %s", pgErr.Message, pgErr.Code, pgErr.Detail)
		}
		return fmt.Errorf("failed to save bet: %w", err)
	}

	if result.RowsAffected() == 0 {
		return FailedOnInsertingData
	}

	return nil
}

func (repo *BetRepository) FindByRound(round *model.Round) ([]*model.Bet, error) {
	query := `
		SELECT 
			b.amount, 
			p.balance, 
			p.status, 
			u.uid, 
			u.name, 
			g.uid, 
			g.action_duration, 
			g.question_duration, 
			r.uid, 
			r.status
		FROM bets b
		JOIN players p ON b.player_id = p.id
		JOIN users u ON p.user_id = u.id
		JOIN lock_stock_games g ON p.game_id = g.id
		JOIN rooms r ON g.room_id = r.id
		WHERE b.round_id = (SELECT id FROM rounds WHERE uid = $1)
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
		var (
			amount           int
			balance          int
			status           model.PlayerStatus
			userUid          string
			username         string
			gameUid          string
			actionDuration   string
			questionDuration string
			roomUid          string
			roomStatus       roomModel.RoomStatus
		)

		err := rows.Scan(
			&amount,
			&balance, &status,
			&userUid, &username,
			&gameUid, &actionDuration, &questionDuration,
			&roomUid, &roomStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bet row: %w", err)
		}

		user := userModel.NewUser(userUid, username)
		room := roomModel.NewRoom(roomUid, roomStatus)
		game := model.NewLockStockGame(gameUid, actionDuration, questionDuration, room)
		player := model.NewPlayer(user, balance, status, game)

		bet := model.NewBet(player, amount, round)
		bets = append(bets, bet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over bet rows: %w", err)
	}

	return bets, nil
}
