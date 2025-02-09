package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"lock-stock-v2/internal/domain/user/model"
	"time"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) FindById(userId string) (*model.User, error) {
	var tempUid string
	var tempName string

	query := `SELECT uid, name FROM users WHERE uid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.db.QueryRow(ctx, query, userId).Scan(&tempUid, &tempName)
	if err != nil {
		return nil, errors.New("user not found: " + err.Error())
	}

	newUser := model.NewUser(tempUid, tempName)
	return newUser, nil
}

func (repo *UserRepository) SaveUser(user *model.User) error {
	query := `
        INSERT INTO users (uid, name)
        VALUES ($1, $2)
        ON CONFLICT (uid) DO UPDATE
        SET name = EXCLUDED.name
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Exec(ctx, query, user.Uid(), user.Name())
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
