package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provee todas las funciones para ejecutar queries de bdd
// y ejecutar transacciones
type Store struct {
	*Queries
	dbPool *pgxpool.Pool
}

// Crea una nueva Store
func NewStore(dbPool *pgxpool.Pool) *Store {
	return &Store{
		dbPool:  dbPool,
		Queries: New(dbPool),
	}
}

// ejecuta una función dentro de una transaccion en bdd
func (store *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.dbPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
