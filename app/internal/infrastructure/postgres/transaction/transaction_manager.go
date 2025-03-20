package transaction

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"lock-stock-v2/external/transaction"
)

type PostgresTransactionManager struct {
	db *pgxpool.Pool
}

func NewPostgresTransactionManager(db *pgxpool.Pool) transaction.TransactionManager {
	return &PostgresTransactionManager{db: db}
}

func (tm *PostgresTransactionManager) Run(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := tm.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction rollback failed: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
