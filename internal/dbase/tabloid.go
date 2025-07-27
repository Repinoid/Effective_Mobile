package dbase

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"emobile/internal/models"
)

// Структура для базы данных.
type DBstruct struct {
	DB *pgxpool.Pool
	//	DB *pgx.Conn
}

func NewPostgresPool(ctx context.Context, DSN string) (*DBstruct, error) {

	poolConfig, err := pgxpool.ParseConfig(DSN)
	//	poolConfig, err := pgxpool.ParseConfig(models.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgxpool config: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbStorage := &DBstruct{}
	dbStorage.DB = pool

	return dbStorage, nil
}

// DataBase PING
func Ping(ctx context.Context) error {
	dataBase, err := NewPostgresPool(ctx, models.DSN)
	if err != nil {
		return err
	}
	defer dataBase.DB.Close()

	err = dataBase.DB.Ping(ctx) // база то открыта ...
	if err != nil {
		models.Logger.Error("No PING ", "error", err.Error())
		return fmt.Errorf("no ping %w", err)
	}
	return nil
}

func (dataBase *DBstruct) AddSub(ctx context.Context, sub models.Subscription) (err error) {

	if sub.End_date == "" {
		order := "INSERT INTO subscriptions(service_name, price, user_id, start_date) VALUES ($1, $2, $3, $4) ;"
		_, err = dataBase.DB.Exec(ctx, order, sub.Service_name, sub.Price, sub.User_id, sub.Sdt)
		return
	}
	order := "INSERT INTO subscriptions(service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) ;"
	_, err = dataBase.DB.Exec(ctx, order, sub.Service_name, sub.Price, sub.User_id, sub.Sdt, sub.Edt)

	return
}

func (dataBase *DBstruct) ReadSub(ctx context.Context, sub models.ReadSubscription) (subs []models.ReadSubscription, err error) {

	return
}
