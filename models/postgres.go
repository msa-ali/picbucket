package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/msa-ali/picbucket/utils"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func DefaultPostgresConfig() PostgresConfig {
	env := utils.GetEnv()
	return PostgresConfig{
		Host:     env.DBHost,
		Port:     env.DBPort,
		User:     env.DBUser,
		Password: env.DBPassword,
		Database: env.DBName,
		SSLMode:  env.DBSSLMode,
	}
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

// Open will open a SQL connection with the provided
// Postgres database. Callers of Open need to ensure
// that the connection is eventually closed via the
// db.Close() method.
func Open(cfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	if dir == "" {
		dir = "."
	}
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
