package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     int
	AppPort    int
}

var Configuration Config

func Load() (*Config, error) {
	// Загружаем .env файл
	// Load will read your env file(s) and load them into ENV for this process.
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("Warning: couldn't load .env file: %v", err)
		// It's important to note that it WILL NOT OVERRIDE an env variable
		// that already exists - consider the .env file to set dev vars or sensible defaults.
		// Не прерываем выполнение, так как переменные могут быть установлены в окружении
		// но и они на самом деле установлены - .env прочитан docker compose
		//
	}

	// Парсим порт приложения
	appPort, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		return nil, err
	}

	// Парсим порт БД
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "postgres"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		AppPort:    appPort,
	}, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
