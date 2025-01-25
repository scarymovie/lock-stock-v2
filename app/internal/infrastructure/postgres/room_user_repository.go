package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	externalDomain "lock-stock-v2/external/domain"
	internalDomain "lock-stock-v2/internal/domain"
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
	ru, ok := roomUser.(*internalDomain.RoomUser)
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

func (repo *RoomUserRepository) FindByRoom(room externalDomain.Room) ([]externalDomain.RoomUser, error) {
	log.Printf("Room type: %T, value: %+v", room, room)

	if room == nil {
		log.Println("FindByRoom called with nil room object!")
		return nil, errors.New("room is nil")
	}

	log.Printf("FindByRoom: about to run query with roomId=%d\n", room.GetRoomId())

	roomID := room.GetRoomId()

	query := `
		SELECT 
			ru.id, 
			r.id, r.uid,
			u.id, u.uid, u.name
		FROM room_users ru
		JOIN rooms r ON ru.room_id = r.id
		JOIN users u ON ru.user_id = u.id
		WHERE ru.room_id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Preparing to query room_users with roomID: %d", roomID)
	rows, err := repo.db.Query(ctx, query, roomID)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, fmt.Errorf("failed to query room_users: %w", err)
	}
	defer rows.Close()
	log.Println("Query executed successfully, processing rows...")

	var results []externalDomain.RoomUser

	for rows.Next() {
		var (
			roomUserID int
			dbRoomID   int
			dbRoomUid  string
			dbUserID   int
			dbUserUid  string
			dbUserName string
		)

		if err := rows.Scan(&roomUserID, &dbRoomID, &dbRoomUid, &dbUserID, &dbUserUid, &dbUserName); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rm := &internalDomain.Room{
			Id:  dbRoomID,
			Uid: dbRoomUid,
		}

		usr := &internalDomain.User{
			Id:   dbUserID,
			Uid:  dbUserUid,
			Name: dbUserName,
		}

		ru := &internalDomain.RoomUser{}
		ru.SetRoom(rm)
		ru.SetUser(usr)

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
