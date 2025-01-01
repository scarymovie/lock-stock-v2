package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	externalDomain "lock-stock-v2/external/domain"
	internalDomain "lock-stock-v2/internal/domain"
	"time"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) FindById(userId string) (externalDomain.User, error) {
	var room internalDomain.User

	query := `SELECT * FROM users WHERE uid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.db.QueryRow(ctx, query, userId).Scan(&room.Id, &room.Uid)
	if err != nil {
		return nil, errors.New("user not found: " + err.Error())
	}

	return &room, nil
}

func (repo *UserRepository) SaveUser(user externalDomain.User) error {

	query := `
		INSERT INTO users (uid)
		VALUES ($1)
		ON CONFLICT (uid) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query, user.GetUserUid())
	if err != nil {
		return fmt.Errorf("failed to save RoomUser: %w", err)
	}

	return nil
}
