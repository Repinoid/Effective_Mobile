package main

import (
	"context"
	"emobile/internal/config"
	"emobile/internal/models"
	"flag"
	"log"
	"log/slog"
	"os"
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

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	models.Logger.Debug("Config", "", cfg)

	if err := Run(ctx); err != nil {
		models.Logger.Error(err.Error())
	}

}

func Run(ctx context.Context) (err error) {

	return
}
