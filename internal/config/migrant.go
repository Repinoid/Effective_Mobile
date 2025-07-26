package config

import (
	//	"emobile/internal/config"
	"emobile/internal/models"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitMigration(cfg Config) (err error) {

	// пока для отладки
	cfg.DBHost = "localhost"
	models.DSN = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	Configuration = cfg

	migrant, err := migrate.New(models.MigrationsPath, models.DSN)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer migrant.Close()

	if err := migrant.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	models.Logger.Debug("migrate", "", migrant)

	version, dirty, err := migrant.Version()
	if err != nil {
		return fmt.Errorf("failed to get migration version: %w", err)
	}
	models.Logger.Debug("Current migration", "version", version, "dirty", dirty)

	return
}
