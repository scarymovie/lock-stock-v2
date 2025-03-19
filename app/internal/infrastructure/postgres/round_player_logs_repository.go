package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"lock-stock-v2/internal/domain/game/model"
	roomModel "lock-stock-v2/internal/domain/room/model"
	userModel "lock-stock-v2/internal/domain/user/model"
	"time"
)

type RoundPlayerLogRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoundPlayerLogRepository(db *pgxpool.Pool) *RoundPlayerLogRepository {
	return &RoundPlayerLogRepository{db: db}
}

func (repo *RoundPlayerLogRepository) FindByRound(round *model.Round) ([]*model.RoundPlayerLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT rpl.player_id, u.uid, u.name, p.balance, p.status, r.id, r.uid, rpl.number, rpl.bets_value, rpl.answer
		FROM round_player_logs rpl
		JOIN rounds r ON rpl.round_id = r.id
		JOIN players p ON rpl.player_id = p.user_id
		JOIN users u ON p.user_id = u.id
		JOIN lock_stock_games g ON r.game_id = g.id
		JOIN rooms rm ON g.room_id = rm.id
		WHERE r.uid = $1
	`

	rows, err := repo.db.Query(ctx, query, round.Uid())
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var logs []*model.RoundPlayerLog
	for rows.Next() {
		var playerID int
		var userID string
		var username string
		var email string
		var balance int
		var status model.PlayerStatus
		var roundID int
		var roundUID string
		var number uint
		var betsValue uint
		var answer *uint
		var gameId string
		var actionDuration string
		var questionDuration string
		var roomId string
		var roomStatus roomModel.RoomStatus

		if err := rows.Scan(
			&playerID, &userID, &username, &email, &balance, &status,
			&roundID, &roundUID, &number, &betsValue, &answer, &gameId, &actionDuration,
			&questionDuration, &roomId, &roomStatus,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		user := userModel.NewUser(userID, username)

		room := roomModel.NewRoom(roomId, roomStatus)

		game := model.NewLockStockGame(gameId, actionDuration, questionDuration, room)

		player := model.NewPlayer(user, balance, status, game)

		log := model.NewRoundPlayerLog(player, round, number, betsValue)
		log.SetAnswer(answer)

		logs = append(logs, log)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while iterating rows: %w", rows.Err())
	}

	return logs, nil
}

func (repo *RoundPlayerLogRepository) FindByRoundAndUser(round *model.Round, user *userModel.User) (*model.RoundPlayerLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println(round.Uid(), user.Uid())

	query := `
		SELECT 
			rpl.player_id, u.uid, u.name, p.balance, p.status, 
			r.id, r.uid, 
			rpl.number, rpl.bets_value, rpl.answer, 
			g.uid, g.action_duration, g.question_duration, 
			rm.id, rm.uid, rm.status
		FROM round_player_logs rpl
		JOIN rounds r ON rpl.round_id = r.id
		JOIN players p ON rpl.player_id = p.user_id
		JOIN users u ON p.user_id = u.id
		JOIN lock_stock_games g ON r.game_id = g.id
		JOIN rooms rm ON g.room_id = rm.id
		WHERE r.uid = $1 AND u.uid = $2
	`

	row := repo.db.QueryRow(ctx, query, round.Uid(), user.Uid())

	var playerID int
	var userID, username string
	var balance int
	var status model.PlayerStatus
	var roundID int
	var roundUID string
	var number, betsValue uint
	var answer *uint
	var gameUID, actionDuration, questionDuration string
	var roomUID string
	var roomStatus roomModel.RoomStatus

	if err := row.Scan(
		&playerID, &userID, &username, &balance, &status,
		&roundID, &roundUID, &number, &betsValue, &answer,
		&gameUID, &actionDuration, &questionDuration,
		&roomUID, &roomStatus,
	); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	room := roomModel.NewRoom(roomUID, roomStatus)

	game := model.NewLockStockGame(gameUID, actionDuration, questionDuration, room)

	player := model.NewPlayer(user, balance, status, game)

	log := model.NewRoundPlayerLog(player, round, number, betsValue)
	log.SetAnswer(answer)

	return log, nil
}

func (repo *RoundPlayerLogRepository) Save(log *model.RoundPlayerLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var playerID, roundID int
	err := repo.db.QueryRow(ctx, `
    SELECT p.id 
    FROM players p 
    JOIN users u ON p.user_id = u.id 
    WHERE u.uid = $1
`, log.Player().User().Uid()).Scan(&playerID)
	if err != nil {
		return fmt.Errorf("failed to fetch player ID: %w", err)
	}

	err = repo.db.QueryRow(ctx, `
    SELECT r.id 
    FROM rounds r 
    WHERE r.uid = $1
`, log.Round().Uid()).Scan(&roundID)
	if err != nil {
		return fmt.Errorf("failed to fetch round ID: %w", err)
	}

	query := `
    INSERT INTO round_player_logs (player_id, round_id, number, bets_value, answer)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (player_id, round_id) 
    DO UPDATE SET number = EXCLUDED.number, bets_value = EXCLUDED.bets_value, answer = EXCLUDED.answer
`
	_, err = repo.db.Exec(ctx, query, playerID, roundID, log.Number(), log.BetsValue(), log.Answer())

	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
