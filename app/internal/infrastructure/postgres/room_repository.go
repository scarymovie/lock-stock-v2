package postgres

import (
	"context"
	"errors"
	"fmt"
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

func (repo *RoomRepository) FindById(roomId string) (externalDomain.Room, error) {
	var room domain.Room

	query := `SELECT id, uid FROM rooms WHERE uid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.db.QueryRow(ctx, query, roomId).Scan(&room.Id, &room.Uid)
	if err != nil {
		return nil, errors.New("room not found: " + err.Error())
	}

	return &room, nil
}

func (repo *RoomRepository) GetAll() ([]externalDomain.Room, error) {
	query := `SELECT id, uid FROM rooms`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []externalDomain.Room
	for rows.Next() {
		var r domain.Room
		if err := rows.Scan(&r.Id, &r.Uid); err != nil {
			return nil, err
		}
		rooms = append(rooms, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (repo *RoomRepository) Save(room externalDomain.Room) error {
	query := `
		INSERT INTO rooms (uid)
		VALUES ($1)
		ON CONFLICT (uid) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query, room.GetRoomUid())
	if err != nil {
		return fmt.Errorf("failed to save room: %w", err)
	}

	return nil
}

func (repo *RoomRepository) UpdateRoomStatus(room externalDomain.Room) error {
	query := `
		UPDATE rooms
		SET status = $1
		WHERE uid = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query, room.GetRoomStatus(), room.GetRoomUid())
	if err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	return nil
}
