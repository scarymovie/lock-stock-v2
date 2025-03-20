package transaction

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type TransactionManager interface {
	Run(ctx context.Context, fn func(tx pgx.Tx) error) error
}
