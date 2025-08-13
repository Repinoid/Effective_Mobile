package handlera

import (
	"emobile/internal/models"
)

type InterStruct struct {
	Inter models.SubscriptionStorage
	//DB    *pgxpool.Pool
	//	DB *pgx.Conn
}

// func NewUserHandler(Inter models.SubscriptionStorage) *DBstruct {
// 	return &DBstruct{Inter: Inter}
// }

// type Handlers interface {
// 	DBPinger(rwr http.ResponseWriter, req *http.Request)
// 	CreateSub(rwr http.ResponseWriter, req *http.Request)
// 	ReadSub(rwr http.ResponseWriter, req *http.Request)
// 	ListSub(rwr http.ResponseWriter, req *http.Request)
// 	UpdateSub(rwr http.ResponseWriter, req *http.Request)
// 	DeleteSub(rwr http.ResponseWriter, req *http.Request)
// 	SumSub(rwr http.ResponseWriter, req *http.Request)
// }

// func NewPostgresPool(ctx context.Context, DSN string) (*DBstruct, error) {

// 	poolConfig, err := pgxpool.ParseConfig(DSN)
// 	//	poolConfig, err := pgxpool.ParseConfig(models.DSN)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse configuration: %w", err)
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create connection pool: %w", err)
// 	}

// 	if err := pool.Ping(ctx); err != nil {
// 		return nil, fmt.Errorf("failed to ping the database: %w", err)
// 	}

// 	// dbStorage := &DBstruct{DB: pool}
// 	// dbStorage.DB = pool

// 	return &DBstruct{DB: pool}, nil
// }
