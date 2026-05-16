package migrate

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/aga-absolut/Vault-System/internal/config"
	"github.com/aga-absolut/Vault-System/internal/interceptors/logger"
	"github.com/pressly/goose/v3"
)

//go:embed migrations
var MigrationsFS embed.FS

func MustInitMigrations(cfg *config.Config, logger *logger.Logger) {
	logger.Info("Starting database migrations...")

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		logger.Errorw("Failed to open database connection for migrations", "error", err)
		panic(fmt.Errorf("sql.Open: %w", err))
	}
	defer db.Close()

	goose.SetBaseFS(MigrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(fmt.Errorf("goose.SetDialect: %w", err))
	}

	if err := goose.Up(db, "migrations"); err != nil {
		logger.Errorw("Migrations failed", "error", err)
		panic(fmt.Errorf("goose.Up: %w", err))
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		logger.Errorw("Failed to get migration version", "error", err)
	} else {
		logger.Infow("Migrations completed successfully", "version", version)
	}
}

func MustDownMigrations(cfg *config.Config, logger *logger.Logger) {
	logger.Info("Starting database rollback...")

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		logger.Errorw("Failed to open database", "error", err)
		panic(err)
	}
	defer db.Close()

	goose.SetBaseFS(MigrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	// Откатываем ОДНУ последнюю миграцию
	if err := goose.Down(db, "migrations"); err != nil {
		logger.Errorw("Migration down failed", "error", err)
		panic(err)
	}

	version, _ := goose.GetDBVersion(db)
	logger.Infow("Rollback completed", "new_version", version)
}
