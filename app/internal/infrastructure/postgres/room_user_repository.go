package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	externalDomain "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain"
	"log"
	"time"
)

type RoomUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoomUserRepository(db *pgxpool.Pool) *RoomUserRepository {
	return &RoomUserRepository{db: db}
}

func (repo *RoomUserRepository) Save(roomUser externalDomain.RoomUser) error {
	ru, ok := roomUser.(*domain.RoomUser)
	if !ok {
		return errors.New("invalid RoomUser type")
	}

	roomId := ru.GetRoom().GetRoomId()
	userId := ru.GetUser().GetUserId()
	log.Println(roomId)
	log.Println(userId)

	query := `
		INSERT INTO room_users (room_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (room_id, user_id) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query, roomId, userId)
	if err != nil {
		return fmt.Errorf("failed to save RoomUser: %w", err)
	}

	return nil
}
