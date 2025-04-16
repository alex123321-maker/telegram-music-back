package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(ctx context.Context, dsn string) error {
	var err error
	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	if err := DB.Ping(ctx); err != nil {
		return fmt.Errorf("ошибка пинга БД: %w", err)
	}

	fmt.Println("🟢 Подключение к БД установлено")
	return nil
}
