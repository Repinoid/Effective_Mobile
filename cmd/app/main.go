package main

import (
	"context"
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

	"github.com/gorilla/mux"
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

	err = initMigration()
	if err != nil {
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/ping", handlera.DBPinger).Methods("GET")

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	hostaport := fmt.Sprintf("%s:%d", models.Config.DBHost, models.Config.AppPort)

	srv := &http.Server{
		Addr:    hostaport,
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Println("Server started on ", hostaport)
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
		models.Logger.Error("Shutdown error", "", err.Error())
	} else {
		log.Println("Server stopped gracefully")
	}

	//err = http.ListenAndServe(fmt.Sprintf("%s:%d", models.Config.DBHost, models.Config.AppPort), router)

	return

}
