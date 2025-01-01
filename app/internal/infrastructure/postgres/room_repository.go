package postgres

import (
	"context"
	"errors"
	externalDomain "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoomRepository(db *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{db: db}
}

// FindById ищет комнату по ID.
func (repo *RoomRepository) FindById(roomId string) (externalDomain.Room, error) {
	var room domain.Room

	query := `SELECT id FROM room WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.db.QueryRow(ctx, query, roomId).Scan(&room.Id)
	if err != nil {
		return nil, errors.New("room not found: " + err.Error())
	}

	return &room, nil
}
