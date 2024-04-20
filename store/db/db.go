package db

import (
	"database/sql"
	"fmt"
	"github.com/dwadp/attendance-api/config"
	_ "github.com/lib/pq"
)

func New(cfg *config.Database) (*sql.DB, error) {
	dsn := BuildDSN(cfg)
	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func BuildDSN(cfg *config.Database) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)
}
