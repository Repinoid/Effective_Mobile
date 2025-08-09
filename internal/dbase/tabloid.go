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
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	dbStorage := &DBstruct{}
	dbStorage.DB = pool

	return dbStorage, nil
}

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

func (dataBase *DBstruct) AddSub(ctx context.Context, sub models.Subscription) (cTag pgconn.CommandTag, err error) {

	// if sub.End_date.(time.Time).IsZero() {
	// 	order := "INSERT INTO subscriptions(service_name, price, user_id, start_date) VALUES ($1, $2, $3, $4) ;"
	// 	cTag, err = dataBase.DB.Exec(ctx, order, sub.Service_name, sub.Price, sub.User_id, sub.Start_date)
	// 	return
	// }
	order := "INSERT INTO subscriptions(service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) ;"
	cTag, err = dataBase.DB.Exec(ctx, order, sub.Service_name, sub.Price, sub.User_id, sub.Start_date, sub.End_date)

	return
}

func (dataBase *DBstruct) ListSub(ctx context.Context, pageSize, offset int) (subs []models.Subscription, err error) {

	order := "SELECT service_name, price, user_id, start_date, end_date FROM subscriptions LIMIT $1 OFFSET $2"
	rows, err := dataBase.DB.Query(ctx, order, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sub := models.Subscription{}
		if err := rows.Scan(&sub.Service_name, &sub.Price, &sub.User_id, &sub.Start_date, &sub.End_date); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return
}

func (dataBase *DBstruct) ReadSub(ctx context.Context, sub models.Subscription) (subs []models.Subscription, err error) {

	order := `
		SELECT service_name, price, user_id, start_date, end_date FROM subscriptions WHERE 
		service_name=$1 AND 
		($2::int = 0 OR price = $2::int) AND
		user_id=$3::uuid AND
		(start_date <= $4 OR $4 = '0001-01-01 00:00:00') AND 
		(end_date >= $5 OR $5 = '0001-01-01 00:00:00' OR end_date ='0001-01-01 00:00:00') ;
	`

	// order := "SELECT service_name, price, user_id, start_date, end_date FROM subscriptions WHERE " +
	// 	"service_name=$1 AND " +
	// 	"($2::int = 0 OR price = $2::int) AND " +
	// 	"user_id=$3::uuid AND " +

	// 	"(start_date <= $4 OR $4 = '0001-01-01 00:00:00') AND " +
	// 	"(end_date >= $5 OR $5 = '0001-01-01 00:00:00' OR end_date ='0001-01-01 00:00:00');"

	rows, err := dataBase.DB.Query(ctx, order, sub.Service_name, sub.Price, sub.User_id, sub.Start_date, sub.End_date)
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
		sub.Start_date = sdt.Time
		sub.End_date = edt.Time
		subs = append(subs, sub)
	}

	return
}

func (dataBase *DBstruct) UpdateSub(ctx context.Context, sub models.Subscription) (cTag pgconn.CommandTag, err error) {

	order := "UPDATE subscriptions SET " +
		"price = CASE WHEN $1::int != 0 THEN $1 ELSE price END, " +
		"start_date=$2, " +
		"end_date=$3 " +
		"WHERE service_name=$4 AND user_id=$5::uuid;"

	cTag, err = dataBase.DB.Exec(ctx, order, sub.Price, sub.Start_date, sub.End_date, sub.Service_name, sub.User_id)

	return
}

func (dataBase *DBstruct) DeleteSub(ctx context.Context, sub models.Subscription) (cTag pgconn.CommandTag, err error) {

	order := "DELETE FROM subscriptions WHERE " +
		"($1 = '' OR service_name = $1) AND " +
		"( ($2::int = 0) OR ($2::int != 0 AND price = $2::int) ) AND " +
		"($3 = '' OR user_id = $3::uuid) AND " +

		"(start_date <= $4 OR $4 = '0001-01-01 00:00:00') AND " +
		"(end_date >= $5 OR $5 = '0001-01-01 00:00:00' OR end_date ='0001-01-01 00:00:00');"

	cTag, err = dataBase.DB.Exec(ctx, order, sub.Service_name, sub.Price, sub.User_id, sub.Start_date, sub.End_date)
	if err != nil {
		models.Logger.Error("Delete", "", err.Error())
	}

	return
}

func (dataBase *DBstruct) SumSub(ctx context.Context, sub models.Subscription) (summa int64, err error) {

	//  если конечная дата подписки не задана - устанавлiваем в максимально возможное значение
	//  			Жаль только — жить в эту пору прекрасную уж не придется — ни мне, ни тебе ©
	if sub.End_date == nil {
		sub.End_date = time.Date(9999, time.December, 31, 23, 59, 59, 999999999, time.UTC)
	}

	// GREATEST($3::DATE, start_date) - начало общего интервала подписка-условие, LEAST($4::DATE, end_date) - окончание
	// разница (конец минус начало) может быть отрицательной (отрезки не пересекаются),
	// поэтому проверка условия dv.effective_start <= dv.effective_end
	order := `
		WITH date_vars AS (
			SELECT id, 
			GREATEST($3::DATE, start_date) AS effective_start, 
			LEAST($4::DATE, end_date) AS effective_end,
			AGE( LEAST($4::DATE, end_date), GREATEST($3::DATE, start_date) ) AS age_interval
			FROM subscriptions
		)
		SELECT SUM(
			s.price *
			(
				-- разница в месяцах ПЛЮС 1, т.к. месяц начала и окончания подписки могут совпадать
				EXTRACT(YEAR FROM age_interval) * 12 +
				EXTRACT(MONTH FROM age_interval) + 1
			)
			) AS total_price
		FROM subscriptions s
		JOIN date_vars dv USING(id)
		-- наименование подписки - либо пусто, либо соответствие табличному
		WHERE ($1 = '' OR s.service_name = $1)
		AND ($2 = '' OR s.user_id = $2::UUID)
		-- условие пересечения временнЫх отрезков
		AND dv.effective_start <= dv.effective_end;
	`

	// order := `
	// SELECT COALESCE(SUM(price *
	// 	((EXTRACT(YEAR FROM LEAST($4::date, end_date)) - EXTRACT(YEAR FROM GREATEST($3::date, start_date))) * 12 +
	// 	(EXTRACT(MONTH FROM LEAST($4::date, end_date)) - EXTRACT(MONTH FROM GREATEST($3::date, start_date))) + 1
	// 		)), 0) AS total_price
	// 	FROM subscriptions
	//  	WHERE ($1 = '' OR service_name = $1)
	//  	AND ($2 = '' OR user_id = $2::uuid)
	//  	AND GREATEST($3, start_date) <= LEAST($4, end_date)
	// `

	row := dataBase.DB.QueryRow(ctx, order, sub.Service_name, sub.User_id, sub.Start_date, sub.End_date)
	summa = 0
	err = row.Scan(&summa)
	if err != nil {
		return 0, err
	}

	return
}

func (dataBase *DBstruct) CloseDB() {
	dataBase.DB.Close()
}

//  docker exec -it pcontB psql -U testuser -d testdb -c "select * from subscriptions"

// docker exec -it pcontB psql -U testuser -d testdb -c "SELECT service_name, start_date, end_date, EXTRACT(MONTH FROM end_date) FROM subscriptions"
