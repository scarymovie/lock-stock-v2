package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	api "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain"
	"time"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) FindById(userId string) (api.User, error) {
	var room domain.User

	query := `SELECT * FROM users WHERE uid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.db.QueryRow(ctx, query, userId).Scan(&room.Id, &room.Uid)
	if err != nil {
		return nil, errors.New("user not found: " + err.Error())
	}

	return &room, nil
}
