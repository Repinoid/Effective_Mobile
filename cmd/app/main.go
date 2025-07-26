package main

import (
	"context"
	"emobile/internal/config"
	"emobile/internal/dbase"
	"emobile/internal/models"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {

	ctx := context.Background()

	// уровень логирования по умолчанию Info
	Level := slog.LevelInfo
	// Если есть флаг -debug
	debugFlag := flag.Bool("debug", false, "установка Минимального уровня логирования DEBUG")
	flag.Parse()
	if *debugFlag {
		Level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     Level,
		AddSource: true, // Добавлять информацию об исходном коде
	})
	models.Logger = slog.New(handler)
	slog.SetDefault(models.Logger)

	models.Logger.Debug("Log", "level", Level)

	if err := Run(ctx); err != nil {
		models.Logger.Error(err.Error())
	}

}

func Run(ctx context.Context) (err error) {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	models.Logger.Debug("Config", "", cfg)

	pool, err := dbase.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf("Failed connect 2 db", err)
	}
	// dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
	// 	cfg.DBUser, cfg.DBPassword, "localhost", cfg.DBPort, cfg.DBName)
	// cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	migrationsPath := "file://../../migrations"

	// 1. Создаём драйвер для существующего подключения
	//	driver, err := postgres.WithInstance(pool, &cfg)

	// Convert pgxpool to stdlib DB
	db := stdlib.OpenDBFromPool(pool)

	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	migrant, err := migrate.NewWithDatabaseInstance(
		migrationsPath, // Полный корректный путь к миграциям
		"postgres",     // Имя драйвера БД
		db,
	)

	// migrant, err := migrate.New(migrationsPath, dsn)
	// if err != nil {
	// 	return fmt.Errorf("failed to create migrate instance: %w", err)
	// }
	// defer migrant.Close()

	// if err := migrant.Steps(-1); err != nil {
	// 	return err
	// }

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
