package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func InitDB(ctx context.Context) (*pgx.Conn, error) {
	var err error

	url, err := getConnectionString()

	if err != nil {
		return nil, err
	}

	conn, err = pgx.Connect(ctx, url)

	if err != nil {
		err = &databaseConnectionError{Err: err}
		log.Printf("%v", err)
		return nil, err
	}

	return conn, err
}

func BeginTx(ctx context.Context) (pgx.Tx, error) {
	return conn.Begin(ctx)
}

func getConnectionString() (string, error) {
	user := os.Getenv("PG_USER")
	if user == "" {
		return "", errors.New("PG_USER environment variable is not set")
	}

	pass := os.Getenv("PG_PASS")
	if pass == "" {
		return "", errors.New("PG_PASS environment variable is not set")
	}

	host := os.Getenv("PG_HOST")
	if host == "" {
		return "", errors.New("PG_HOST environment variable is not set")
	}

	port := os.Getenv("PG_PORT")
	if port == "" {
		return "", errors.New("PG_PORT environment variable is not set")
	}

	dbName := os.Getenv("PG_NAME")
	if dbName == "" {
		return "", errors.New("PG_NAME environment variable is not set")
	}

	sslMode := os.Getenv("PG_SSL_MODE")

	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbName, sslMode), nil
}
