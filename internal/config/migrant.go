package config

import (
	//	"emobile/internal/config"
	"context"
	"emobile/internal/models"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitMigration(ctx context.Context, cfg Config) (err error) {

	// если сервер запущен в контейнере, в нём есть переменная окружения MIGRATIONS_PATH
	enva, exists := os.LookupEnv("MIGRATIONS_PATH")
	if exists {
		// задаём путь к файлам миграции в самом контейнере
		models.MigrationsPath = enva
	} else {
		// если нет "MIGRATIONS_PATH" значит приложение запущено не в контейнере и хост localhost
		cfg.DBHost = "localhost"
	}

	models.DSN = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	Configuration = cfg

	// PING data base check
	err = CheckBase(ctx, models.DSN)
	if err != nil {
		models.Logger.Error("No", "DB", err)
		return
	}

	models.Logger.Info("DB ok", "DSN", models.DSN)

	migrant, err := migrate.New(models.MigrationsPath, models.DSN)
	if err != nil {
		models.Logger.Error("migrate", "MigrationsPath", models.MigrationsPath, "DSN", models.DSN, "ERR", err)
		pureFile, ok := strings.CutPrefix(models.MigrationsPath, "file://")
		if !ok {
			models.Logger.Error("no prefix file://")
		}
		fileInfo, errf := os.Stat(pureFile)
		_ = fileInfo
		pwd, _ := os.Getwd()
		if errf != nil {
			models.Logger.Error("no file ", "", pureFile, "err", errf, "pwd", pwd)
		} else {
			models.Logger.Debug("ok", "exist", pureFile)
		}
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

func getDir(path string) (dirac []string, err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		dirac = append(dirac, file.Name())
	}
	return
}

func CheckBase(ctx context.Context, DSN string) (err error) {

	poolConfig, err := pgxpool.ParseConfig(DSN)
	//	poolConfig, err := pgxpool.ParseConfig(models.DSN)
	if err != nil {
		models.Logger.Error("No", "ParseConfig", err)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		models.Logger.Error("No", "pgxpool.NewWithConfig", err, "PoolConfig", poolConfig)
		return
	}

	if err = pool.Ping(ctx); err != nil {
		models.Logger.Error("No", "Ping", err)
		return
	}
	pool.Close()

	return
}
