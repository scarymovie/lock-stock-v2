package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/internal/domain/room/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoomRepository(db *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{db: db}
}

func (repo *RoomRepository) FindById(roomId string) (*model.Room, error) {
	var tempUid string
	var tempStatus model.RoomStatus

	query := `SELECT uid, status FROM rooms WHERE uid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.db.QueryRow(ctx, query, roomId).Scan(&tempUid, &tempStatus)
	if err != nil {
		return nil, errors.New("room not found: " + err.Error())
	}

	return model.NewRoom(tempUid, tempStatus), nil
}

func (repo *RoomRepository) GetPending() ([]*model.Room, error) {
	var tempId int
	var tempUid string
	var tempStatus model.RoomStatus

	query := `SELECT * FROM rooms where rooms.status = 'pending'`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*model.Room
	for rows.Next() {
		if err := rows.Scan(&tempId, &tempUid, &tempStatus); err != nil {
			return nil, err
		}
		room := model.NewRoom(tempUid, tempStatus)
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (repo *RoomRepository) Save(room *model.Room) error {
	query := `
		INSERT INTO rooms (uid)
		VALUES ($1)
		ON CONFLICT (uid) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query, room.Uid())
	if err != nil {
		return fmt.Errorf("failed to save room: %w", err)
	}

	return nil
}

func (repo *RoomRepository) UpdateRoomStatus(ctx context.Context, tx pgx.Tx, room *model.Room) error {
	query := `
		UPDATE rooms
		SET status = $1
		WHERE uid = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tx.Exec(ctx, query, room.Status(), room.Uid())
	if err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	return nil
}
