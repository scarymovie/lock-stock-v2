package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	roomModel "lock-stock-v2/internal/domain/room/model"
	userModel "lock-stock-v2/internal/domain/room/model"
	roomUserModel "lock-stock-v2/internal/domain/room_user/model"
	UserModel "lock-stock-v2/internal/domain/user/model"
	"log"
	"time"
)

type RoomUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoomUserRepository(db *pgxpool.Pool) *RoomUserRepository {
	return &RoomUserRepository{db: db}
}

func (repo *RoomUserRepository) Save(roomUser *roomUserModel.RoomUser) error {
	roomUid := roomUser.Room().Uid()
	userUid := roomUser.User().Uid()

	query := `
		INSERT INTO room_users (room_id, user_id)
		VALUES ((select id from rooms where uid = $1), (select id from users where uid = $2))
		ON CONFLICT (room_id, user_id) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query, roomUid, userUid)
	if err != nil {
		return fmt.Errorf("failed to save RoomUser: %w", err)
	}

	return nil
}

func (repo *RoomUserRepository) FindByRoom(room *roomModel.Room) ([]*roomUserModel.RoomUser, error) {
	query := `
		SELECT r.uid, r.status, u.uid, u.name FROM room_users ru
		JOIN rooms r ON ru.room_id = r.id
		JOIN users u ON ru.user_id = u.id
		WHERE ru.room_id = (select id from rooms where uid = $1)
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := repo.db.Query(ctx, query, room.Uid())
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, fmt.Errorf("failed to query room_users: %w", err)
	}
	defer rows.Close()
	log.Println("Query executed successfully, processing rows...")

	var results []*roomUserModel.RoomUser

	for rows.Next() {
		var (
			roomUid    string
			roomStatus userModel.RoomStatus
			userUid    string
			userName   string
		)

		if err := rows.Scan(&roomUid, &roomStatus, &userUid, &userName); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rm := roomModel.NewRoom(userUid, roomStatus)
		usr := UserModel.NewUser(userUid, userName)
		ru := roomUserModel.NewRoomUser(rm, usr)

		results = append(results, ru)
	}

	if err := rows.Err(); err != nil {
		log.Printf("FindByRoom: iteration error: %v\n", err)
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	if len(results) == 0 {
		log.Println("FindByRoom: no users found in this room.")
		return nil, errors.New("no users found in this room")
	}

	log.Printf("FindByRoom: found %d room-users\n", len(results))
	return results, nil
}
