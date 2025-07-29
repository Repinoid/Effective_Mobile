package config

import (
	//	"emobile/internal/config"
	"emobile/internal/models"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitMigration(cfg Config) (err error) {

	// пока для отладки
	//	cfg.DBHost = "localhost"
	enva, exists := os.LookupEnv("BASE_HOST")
	if exists {
		cfg.DBHost = enva
	}
	models.DSN = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	Configuration = cfg

	// если сервер запущен в контейнере, в нём есть переменная окружения MIGRATIONS_PATH
	enva, exists = os.LookupEnv("MIGRATIONS_PATH")
	if exists {
		models.MigrationsPath = enva
	}

	migrant, err := migrate.New(models.MigrationsPath, models.DSN)
	if err != nil {
		models.Logger.Error("migrate", "MigrationsPath", models.MigrationsPath, "DSN", models.DSN, "", err)
		pureFile, ok := strings.CutPrefix(models.MigrationsPath, "file://")
		if !ok {
			models.Logger.Error("no prefix file://")
		}
		fileInfo, errf := os.Stat(pureFile)
		_ = fileInfo
		if errf != nil {
			models.Logger.Error("no file", "", pureFile, "err", errf)
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
