package models

import "log/slog"

var (
	Logger *slog.Logger
	MigrationsPath = "file://../../migrations"
	DSN = ""
)
