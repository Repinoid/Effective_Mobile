package models

import (
	"emobile/internal/config"
	"log/slog"
)

var (
	Logger         *slog.Logger
	MigrationsPath = "file://../../migrations"
	DSN            = ""
	Config         config.Config
)
