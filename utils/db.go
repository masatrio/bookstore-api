package utils

import (
	"context"
	"database/sql"
	"log"
)

const (
	TransactionContextKey = "DB-TRX"
)

// ExecContextWithPreparedReturningID executes a prepared statement that returns a single ID using
// a transaction from the context if available, or the database otherwise.
func ExecContextWithPreparedReturningID(ctx context.Context, db *sql.DB, query string, args ...interface{}) (int64, error) {
	// Check if there's an active transaction in the context
	tx, ok := ctx.Value(TransactionContextKey).(*sql.Tx)
	var stmt *sql.Stmt
	var err error

	if ok {
		// Prepare statement on the transaction
		stmt, err = tx.PrepareContext(ctx, query)
		if err != nil {
			log.Printf("Transaction preparation error: %v", err)
			return 0, err
		}
		defer stmt.Close()

		var id int64
		err = stmt.QueryRowContext(ctx, args...).Scan(&id)
		if err != nil {
			log.Printf("Transaction query execution error: %v", err)
			return 0, err
		}
		return id, nil
	}

	// Prepare the statement on the regular DB connection
	stmt, err = db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("DB preparation error: %v", err)
		return 0, err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRowContext(ctx, args...).Scan(&id)
	if err != nil {
		log.Printf("DB query execution error: %v", err)
		return 0, err
	}

	return id, nil
}

// prepareAndQueryRowContext prepares a statement with transaction support and executes QueryRowContext.
func PrepareAndQueryRowContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) *sql.Row {
	tx, ok := ctx.Value(TransactionContextKey).(*sql.Tx)
	if ok {
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return nil
		}
		defer stmt.Close()
		return stmt.QueryRowContext(ctx, args...)
	}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}
	defer stmt.Close()

	return stmt.QueryRowContext(ctx, args...)
}

// prepareAndQueryContext prepares a statement with transaction support and executes QueryContext.
func PrepareAndQueryContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	tx, ok := ctx.Value(TransactionContextKey).(*sql.Tx)
	if ok {
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		return stmt.QueryContext(ctx, args...)
	}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, args...)
}
