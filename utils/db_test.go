package utils

import (
	"context"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestExecContextWithPreparedReturningID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	query := "INSERT INTO test_table (name) VALUES (?) RETURNING id"
	expectedID := int64(1)

	log.Printf("Preparing SQL: %s", query)

	mock.ExpectPrepare("INSERT INTO test_table \\(name\\) VALUES \\(\\?\\) RETURNING id").
		ExpectQuery().
		WithArgs("test-name").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	id, err := ExecContextWithPreparedReturningID(ctx, db, query, "test-name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, expectedID, id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctxWithTx := context.WithValue(ctx, TransactionContextKey, tx)

	log.Printf("Preparing SQL for transaction: %s", query)

	mock.ExpectPrepare("INSERT INTO test_table \\(name\\) VALUES \\(\\?\\) RETURNING id").
		ExpectQuery().
		WithArgs("test-name").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	id, err = ExecContextWithPreparedReturningID(ctxWithTx, db, query, "test-name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, expectedID, id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPrepareAndQueryRowContext(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	query := "SELECT id FROM test_table WHERE name = ?"
	expectedID := 1

	ctx := context.Background()

	mock.ExpectPrepare(query).
		ExpectQuery().
		WithArgs("test-name").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	row := PrepareAndQueryRowContext(ctx, db, query, "test-name")

	var id int
	err = row.Scan(&id)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctxWithTx := context.WithValue(ctx, TransactionContextKey, tx)

	mock.ExpectPrepare(query).
		ExpectQuery().
		WithArgs("test-name").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	row = PrepareAndQueryRowContext(ctxWithTx, db, query, "test-name")

	err = row.Scan(&id)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPrepareAndQueryContext(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	query := "SELECT id FROM test_table WHERE name = ?"
	expectedID := 1

	ctx := context.Background()

	mock.ExpectPrepare(query).
		ExpectQuery().
		WithArgs("test-name").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	rows, err := PrepareAndQueryContext(ctx, db, query, "test-name")
	assert.NoError(t, err)
	defer rows.Close()

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		assert.NoError(t, err)
	}
	assert.Equal(t, expectedID, id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctxWithTx := context.WithValue(ctx, TransactionContextKey, tx)

	mock.ExpectPrepare(query).
		ExpectQuery().
		WithArgs("test-name").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	rows, err = PrepareAndQueryContext(ctxWithTx, db, query, "test-name")
	assert.NoError(t, err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		assert.NoError(t, err)
	}
	assert.Equal(t, expectedID, id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
