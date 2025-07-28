package main

import (
	"context"
	"emobile/internal/config"
	"emobile/internal/handlera"
	"emobile/internal/models"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "emobile/docs" // docs генерируется Swag CLI

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Subscription API
// @version 1.0
// @description API для управления подписками и проверки состояния БД
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
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
	models.Logger.Debug("Config", "", *cfg)

	err = config.InitMigration(*cfg)
	if err != nil {
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handlera.DBPinger).Methods("GET")
	router.HandleFunc("/add", handlera.CreateSub).Methods("POST")
	router.HandleFunc("/read", handlera.ReadSub).Methods("POST")
	router.HandleFunc("/list", handlera.ListSub).Methods("GET")
	router.HandleFunc("/update", handlera.UpdateSub).Methods("PUT")
	router.HandleFunc("/delete", handlera.DeleteSub).Methods("DELETE")

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	hostaport := fmt.Sprintf("%s:%d", config.Configuration.DBHost, config.Configuration.AppPort)

	srv := &http.Server{
		Addr:    hostaport,
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		models.Logger.Info("Server started", "on", hostaport)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидаем SIGINT (Ctrl+C) или SIGTERM
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-exit
		cancel() // При получении сигнала отменяем контекст
	}()

	// Блокируемся, пока контекст не отменён
	<-ctx.Done()

	// Graceful shutdown с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		models.Logger.Error("Shutdown", "error", err.Error())
	} else {
		models.Logger.Info("Server stopped gracefully")
	}

	//err = http.ListenAndServe(fmt.Sprintf("%s:%d", models.Config.DBHost, models.Config.AppPort), router)

	return

}
