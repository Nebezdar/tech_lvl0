package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(connString string) (*Database, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	return &Database{pool: pool}, nil
}

func (db *Database) SaveOrder(orderID string) error {
	_, err := db.pool.Exec(context.Background(), "INSERT INTO orders (order_id) VALUES ($1) ON CONFLICT (order_id) DO NOTHING", orderID)
	return err
}

func (db *Database) GetOrder(orderID string) (string, error) {
	var result string
	err := db.pool.QueryRow(context.Background(), "SELECT order_id FROM orders WHERE order_id = $1", orderID).Scan(&result)
	return result, err
}
