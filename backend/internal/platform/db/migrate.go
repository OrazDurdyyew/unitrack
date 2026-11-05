package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
)

// migrations — SQL-файлы миграций, вшитые в бинарник.
//
//go:embed migrations/*.sql
var migrations embed.FS

// Migrate применяет все неприменённые миграции. Повторный вызов безопасен.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer func() { _ = sqlDB.Close() }() // ошибку закрытия обёртки нечем обработать: пул живёт дальше

	files, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("каталог миграций: %w", err)
	}
	provider, err := goose.NewProvider(database.DialectPostgres, sqlDB, files, goose.WithVerbose(false))
	if err != nil {
		return fmt.Errorf("инициализация миграций: %w", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("применение миграций: %w", err)
	}
	return nil
}
