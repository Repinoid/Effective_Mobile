package main

import (
	"context"
	"emobile/internal/config"
	"emobile/internal/handlera"
	"emobile/internal/middlas"
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
	//	_ "github.com/swaggo/http-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Subscription Service API
// @version 1.0
// @description API для управления подписками

// main godoc
// @Summary Запуск приложения
// @Description Основная функция запуска сервиса подписок
// @Produce json
// @Param debug query boolean false "Включить debug-логирование" default(false)
// @Success 200 {string} string "Сервис запущен"
// @Failure 500 {string} string "Ошибка сервера"
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

// Run godoc
// @Summary Запуск сервера API
// @Description Инициализирует конфигурацию и запускает HTTP-сервер с роутингом
// @Accept json
// @Produce json
// @Param ctx query string false "Контекст выполнения"
// @Success 200 {string} string "Сервер запущен"
// @Failure 500 {string} string "Ошибка сервера"
func Run(ctx context.Context) (err error) {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	models.Logger.Debug("Config", "", *cfg)

	err = config.InitMigration(ctx, *cfg)
	if err != nil {
		models.Logger.Debug("sleep ...", "", *cfg)
		time.Sleep(900 * time.Second)
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handlera.DBPinger).Methods("GET")
	router.HandleFunc("/add", handlera.CreateSub).Methods("POST")
	router.HandleFunc("/read", handlera.ReadSub).Methods("POST")
	router.HandleFunc("/list", handlera.ListSub).Methods("GET")
	router.HandleFunc("/update", handlera.UpdateSub).Methods("PUT")
	router.HandleFunc("/delete", handlera.DeleteSub).Methods("DELETE")
	router.HandleFunc("/summa", handlera.SumSub).Methods("POST")

	// подключаем middleware логирования
	router.Use(middlas.WithHTTPLogging)
	router.Use(middlas.ErrorLoggerMiddleware)

	router.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		data, _ := os.ReadFile("./docs/swagger.json")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	// Обработчик для Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/swagger.json"), // Указываем путь к JSON
		httpSwagger.DocExpansion("none"),         // Опционально: схлопывать документацию
	))

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//hostaport := fmt.Sprintf("%s:%d", config.Configuration.DBHost, config.Configuration.AppPort)

	// Используем AppHost (или 0.0.0.0) и AppPort для HTTP-сервера
	serverAddr := fmt.Sprintf("%s:%d", config.Configuration.AppHost, config.Configuration.AppPort)

	// srv := &http.Server{
	// 	Addr:    hostaport,
	// 	Handler: router,
	// }

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		fmt.Printf("\nServer started on %s\n\n", serverAddr)
		models.Logger.Info("Server started", "on", serverAddr)
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
