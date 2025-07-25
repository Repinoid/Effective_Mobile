package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     int
	AppPort    int
}

func Load() (*Config, error) {
	port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		AppPort:    port,
	}, nil
}
