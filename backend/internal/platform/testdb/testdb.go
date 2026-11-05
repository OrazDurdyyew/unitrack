// Package testdb поднимает настоящую PostgreSQL для тестов.
//
// По умолчанию — контейнер postgres:17-alpine (переопределяется TESTDB_IMAGE) через
// testcontainers, один на процесс тестов.
// Если задана переменная TESTDB_URL (например, в CI с сервисом postgres), контейнер не
// создаётся, используется указанный сервер. В обоих случаях каждый вызов New получает
// отдельную базу данных с применёнными миграциями, поэтому тесты изолированы друг от друга.
package testdb

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/OrazDurdyyew/unitrack/backend/internal/platform/db"
)

var (
	once     sync.Once
	adminURL string
	initErr  error
)

// New возвращает пул к свежей базе данных с применёнными миграциями.
// База удаляется в t.Cleanup. Падает через t.Fatal, если PostgreSQL недоступна.
func New(t testing.TB) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	once.Do(func() { adminURL, initErr = startServer(ctx) })
	if initErr != nil {
		t.Fatalf("testdb: запуск PostgreSQL: %v", initErr)
	}

	name := "test_" + randomSuffix()
	admin, err := pgxpool.New(ctx, adminURL)
	if err != nil {
		t.Fatalf("testdb: административное подключение: %v", err)
	}
	if _, err := admin.Exec(ctx, `CREATE DATABASE `+name); err != nil {
		admin.Close()
		t.Fatalf("testdb: создание базы %s: %v", name, err)
	}

	pool, err := db.NewPool(ctx, withDatabase(adminURL, name))
	if err != nil {
		admin.Close()
		t.Fatalf("testdb: подключение к %s: %v", name, err)
	}
	if err := db.Migrate(ctx, pool); err != nil {
		pool.Close()
		admin.Close()
		t.Fatalf("testdb: миграции: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
		dropCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = admin.Exec(dropCtx, `DROP DATABASE IF EXISTS `+name+` WITH (FORCE)`)
		admin.Close()
	})
	return pool
}

// startServer возвращает строку подключения к серверу PostgreSQL с правами на CREATE DATABASE.
func startServer(ctx context.Context) (string, error) {
	if u := os.Getenv("TESTDB_URL"); u != "" {
		return u, nil
	}
	startCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	image := os.Getenv("TESTDB_IMAGE")
	if image == "" {
		image = "postgres:17-alpine"
	}
	container, err := postgres.Run(startCtx, image,
		postgres.WithDatabase("unitrack_test"),
		postgres.WithUsername("unitrack"),
		postgres.WithPassword("unitrack"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(90*time.Second),
		),
	)
	if err != nil {
		return "", fmt.Errorf("контейнер postgres: %w", err)
	}
	u, err := container.ConnectionString(startCtx, "sslmode=disable")
	if err != nil {
		return "", fmt.Errorf("строка подключения контейнера: %w", err)
	}
	return u, nil
}

func withDatabase(connURL, dbName string) string {
	u, err := url.Parse(connURL)
	if err != nil {
		return connURL
	}
	u.Path = "/" + dbName
	return u.String()
}

func randomSuffix() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
