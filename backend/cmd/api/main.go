// Команда api — HTTP-сервер UniTrack.
//
// Запуск: переменные окружения UNITRACK_* (см. internal/platform/config).
// Запуск с флагом -healthcheck делает запрос к собственному /api/v1/ready и завершается
// с кодом 0 или 1 — так docker compose проверяет здоровье контейнера без curl.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/OrazDurdyyew/unitrack/backend/internal/health"
	"github.com/OrazDurdyyew/unitrack/backend/internal/platform/config"
	"github.com/OrazDurdyyew/unitrack/backend/internal/platform/db"
	"github.com/OrazDurdyyew/unitrack/backend/internal/platform/httpserver"
)

// version подставляется при сборке: -ldflags "-X main.version=..."
var version = "dev"

const (
	exitOK       = 0
	exitRuntime  = 1
	exitConfig   = 2
	shutdownWait = 10 * time.Second
)

func main() {
	healthcheck := flag.Bool("healthcheck", false, "проверить /api/v1/ready запущенного сервера и выйти")
	flag.Parse()

	if *healthcheck {
		os.Exit(runHealthcheck(os.Getenv))
	}
	os.Exit(run(context.Background(), os.Getenv))
}

func run(ctx context.Context, getenv func(string) string) int {
	cfg, err := config.Load(getenv)
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error("сервер не запущен", "err", err)
		return exitConfig
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(log)
	log.Info("запуск", "version", version, "env", cfg.Env, "addr", cfg.HTTPAddr)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("подключение к базе данных", "err", err)
		return exitRuntime
	}
	defer pool.Close()

	if migrateOnStart(getenv) {
		if err := db.Migrate(ctx, pool); err != nil {
			log.Error("миграции", "err", err)
			return exitRuntime
		}
		log.Info("миграции применены")
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpserver.NewRouter(log, health.New(version, pool)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("слушаю", "addr", cfg.HTTPAddr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP-сервер", "err", err)
			return exitRuntime
		}
	case <-ctx.Done():
		log.Info("получен сигнал остановки")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownWait)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("остановка HTTP-сервера", "err", err)
		return exitRuntime
	}
	log.Info("остановлен")
	return exitOK
}

// migrateOnStart — UNITRACK_MIGRATE_ON_START, по умолчанию true.
func migrateOnStart(getenv func(string) string) bool {
	raw := getenv("UNITRACK_MIGRATE_ON_START")
	if raw == "" {
		return true
	}
	v, err := strconv.ParseBool(raw)
	return err != nil || v
}

// runHealthcheck обращается к собственному /ready. Адрес берётся из UNITRACK_HTTP_ADDR.
func runHealthcheck(getenv func(string) string) int {
	addr := getenv("UNITRACK_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck: некорректный UNITRACK_HTTP_ADDR:", err)
		return exitConfig
	}
	if host == "" {
		host = "127.0.0.1"
	}
	url := "http://" + net.JoinHostPort(host, port) + httpserver.BasePath + "/ready"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		return exitRuntime
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		return exitRuntime
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "healthcheck: статус", res.StatusCode)
		return exitRuntime
	}
	return exitOK
}
