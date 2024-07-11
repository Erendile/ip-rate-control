package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"ip-rate-control/pkg/config"
)

func NewPostgresDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")

	err = createTableIfNotExists(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTableIfNotExists(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ip_requests (
		ip_address VARCHAR(45) PRIMARY KEY,
		first_request TIMESTAMPTZ NOT NULL
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
