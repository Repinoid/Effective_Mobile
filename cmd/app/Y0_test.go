package main

// Basic imports
import (
	"context"
	"emobile/internal/config"
	"emobile/internal/models"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type TS struct {
	suite.Suite
	t   time.Time
	ctx context.Context
}

func (suite *TS) SetupTest() {
	suite.ctx = context.Background()
	suite.t = time.Now()

	models.MigrationsPath = "file://../../migrations"
	err := godotenv.Load("../../.env")
	suite.Require().NoError(err, "No .ENV file load")

	cfg := config.Config{
		DBUser:     config.GetEnv("DB_USER", "postgres"),
		DBPassword: config.GetEnv("DB_PASSWORD", ""),
		DBName:     config.GetEnv("DB_NAME", "postgres"),
		DBHost:     "localhost",
		DBPort:     5432,
		// AppPort:    8080,
		// AppHost:    "localhost",
	}

	models.DSN = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// PING data base check
	err = config.CheckBase(suite.ctx, models.DSN)
	suite.Require().NoError(err, "No DataBase connection")
	err = config.InitMigration(suite.ctx, cfg)
	suite.Require().NoError(err, "Миграция не прошла")

	// delete все записи - маска с пустой структурой models.Subscription{}
	httpc := resty.New().SetBaseURL("http://localhost:8080")
	req := httpc.R().SetHeader("Content-Type", "application/json").SetDoNotParseResponse(false).
		SetBody(models.Subscription{})
	_ = req
	// раскомментировать если надо обунулять базу перед тестами
	// _, err = req.Delete("/delete")
	// suite.Require().NoError(err, "DROP")

}

func (suite *TS) BeforeTest(suiteName, testName string) {
	log.Println("BeforeTest()", suiteName, testName)
}

func (suite *TS) AfterTest(suiteName, testName string) {
	log.Println("AfterTest()", suiteName, testName)
}

func TestExampleTestSuite(t *testing.T) {
	log.Println("before run")

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true, // Добавлять информацию об исходном коде
	})
	models.Logger = slog.New(handler)
	slog.SetDefault(models.Logger)

	suite.Run(t, new(TS))
}
