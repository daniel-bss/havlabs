package db

import (
	"context"
	"fmt"
)

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("TX error: %v, Rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
