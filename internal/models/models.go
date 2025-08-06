package models

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	Logger         *slog.Logger
	MigrationsPath = "file://migrations"
	// MigrationsPath = "file://../../migrations"
	EnvPath = "./.env"
	DSN     = ""
)

type Subscription struct {
	Service_name string `json:"service_name"`       // “Yandex Plus”,
	Price        int64  `json:"price"`              // “price”: 400,
	User_id      string `json:"user_id"`            // “user_id”: “60601fee-2bf1-4721-ae6f-7636e79a0cba”,
	Start_date   string `json:"start_date"`         // “start_date”: “07-2025”
	End_date     string `json:"end_date,omitempty"` // “start_date”: “07-2025”
	Sdt          any    `json:"-"`                  // немаршалемое
	Edt          any    `json:"-"`
	// Sdt          time.Time `json:"-"`                  // немаршалемое
	// Edt          time.Time `json:"-"`
}

type RetStruct struct {
	Name string
	Cunt int64
}

var Inter Interferon

type Interferon interface {
	AddSub(ctx context.Context, sub Subscription) (cTag pgconn.CommandTag, err error)
	ListSub(ctx context.Context, pageSize, offset int) (subs []Subscription, err error)
	ReadSub(ctx context.Context, sub Subscription) (subs []Subscription, err error)
	UpdateSub(ctx context.Context, sub Subscription) (cTag pgconn.CommandTag, err error)
	DeleteSub(ctx context.Context, sub Subscription) (cTag pgconn.CommandTag, err error)
	SumSub(ctx context.Context, sub Subscription) (summa int64, err error)
	CloseDB()
}
