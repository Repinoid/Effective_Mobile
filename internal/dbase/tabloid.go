package dbase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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

func (dataBase *DBstruct) ListSub(ctx context.Context) (subs []models.Subscription, err error) {

	order := "SELECT service_name, price, user_id, start_date, end_date FROM subscriptions"
	rows, err := dataBase.DB.Query(ctx, order)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sub := models.Subscription{}
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

func (dataBase *DBstruct) ReadSub(ctx context.Context, sub models.Subscription) (subs []models.Subscription, err error) {

	// sub.Sdt nilEdt тип time.Time.
	// происходит полная муть если это передавать в Query из-за того что у них нет обычного nil,
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

	// Так как start_date и end_date могут и не присутствовать в запросе, передаём их в order по COALESCE
	order := "SELECT service_name, price, user_id, start_date, end_date FROM subscriptions WHERE " +
		"service_name=$1 AND " +
		"($2::int = 0 OR price = $2::int) AND " +
		"user_id=$3 AND " +

		"(start_date <= $4 OR $4 IS NULL) AND " +
		"(end_date >= $5 OR $5 IS NULL OR end_date IS NULL);"
		// "start_date <= COALESCE($4, start_date) AND " +
		// "end_date >= COALESCE($5, end_date);"

	rows, err := dataBase.DB.Query(ctx, order, sub.Service_name, sub.Price, sub.User_id, nilSdt, nilEdt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sub := models.Subscription{}
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

func (dataBase *DBstruct) UpdateSub(ctx context.Context, sub models.Subscription) (cTag pgconn.CommandTag, err error) {

	// comments on ReadSub
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

	// comments on ReadSub
	order := "UPDATE subscriptions SET " +
		"price = CASE WHEN $1::int != 0 THEN $1 ELSE price END, " +
		//"price=COALESCE($1, price), " +
		"start_date=COALESCE($2, start_date), " +
		"end_date=COALESCE($3, end_date) " +
		"WHERE service_name=$4 AND user_id=$5;"

	cTag, err = dataBase.DB.Exec(ctx, order, sub.Price, nilSdt, nilEdt, sub.Service_name, sub.User_id)

	return
}

func (dataBase *DBstruct) DeleteSub(ctx context.Context, sub models.Subscription) (cTag pgconn.CommandTag, err error) {

	// comments on ReadSub
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

	order := "DELETE FROM subscriptions WHERE " +
		"($1 = '' OR service_name = $1) AND " +
		"( ($2::int = 0 AND price != 0) OR ($2::int != 0 AND price = $2::int) ) AND " +
		"($3 = '' OR user_id = $3) AND " +

		"(start_date <= $4 OR $4 IS NULL) AND " +
		"(end_date >= $5 OR $5 IS NULL OR end_date IS NULL);"

	cTag, err = dataBase.DB.Exec(ctx, order, sub.Service_name, sub.Price, sub.User_id, nilSdt, nilEdt)
	if err != nil {
		models.Logger.Error("Delete", "", err.Error())
	}

	return
}
