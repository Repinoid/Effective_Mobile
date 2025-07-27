package dbase

import (
	"context"
	"database/sql"
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

func (dataBase *DBstruct) ListSub(ctx context.Context) (subs []models.ReadSubscription, err error) {

	order := "SELECT service_name, price, user_id, start_date, end_date FROM subscriptions"
	rows, err := dataBase.DB.Query(ctx, order)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sub := models.ReadSubscription{}
		// Start_time ||& End_time могут быть NULL. поэтому в Scan подставляем переменные sql.NullTime
		var sdt, edt sql.NullTime
		// err := row.Scan(&createdAt)
		// if createdAt.Valid {}
		if err := rows.Scan(&sub.Service_name, &sub.Price, &sub.User_id, &sdt, &edt); err != nil {
			return nil, err
		}
		sub.Sdt = sdt.Time
		sub.Edt = edt.Time
		subs = append(subs, sub)
	}

	return
}

func (dataBase *DBstruct) ReadSub(ctx context.Context, sub models.ReadSubscription) (subs []models.ReadSubscription, err error) {

	// sub.Sdt тип time.Time. происходит полная муть если это передавать в Query из-за того что у них нет номального nil,
	// определяем нулёвость по .IsZero() & прописываем в интерфейс, который и подсовываем в Query
	var nilSdt, nilEdt any
	if sub.Sdt.IsZero() {
		nilSdt = nil
	} else {
		nilSdt = sub.Sdt
	}
	if sub.Edt.IsZero() {
		nilEdt = nil
	} else {
		nilEdt = sub.Edt
	}

	// Так как price, start_date и end_date могут и не присутствовать в запросе, передаём их в order по COALESCE
	// COALESCE возвращает первый ненулевой параметр. Например -
	// COALESCE($2, price) - если price $2 не нуль, возвращается его значение
	// и  получается обычное сравнение price = $2
	// если $2=0 т.е. требование по price в запрос не передано, получается price = price
	//  - всегда TRUE и этот пункт WHERE попросту игнорируется
	order := "SELECT service_name, price, user_id, start_date, end_date FROM subscriptions WHERE " +
		//	order := "SELECT service_name, price, user_id, start_date FROM subscriptions WHERE " +
		//	order := "SELECT start_date FROM subscriptions WHERE " +
		"service_name=$1 AND " +
		"(price = COALESCE($2, price)) AND " +
		"user_id=$3 AND " +
		"start_date <= COALESCE($4, start_date) AND " +
		"end_date >= COALESCE($5, end_date)"

	rows, err := dataBase.DB.Query(ctx, order, sub.Service_name, sub.Price, sub.User_id, nilSdt, nilEdt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sub := models.ReadSubscription{}
		var sdt, edt sql.NullTime
		if err := rows.Scan(&sub.Service_name, &sub.Price, &sub.User_id, &sdt, &edt); err != nil {
			return nil, err
		}
		sub.Sdt = sdt.Time
		sub.Edt = edt.Time
		subs = append(subs, sub)
	}

	return
}
