// Package config читает настройки сервера из переменных окружения.
//
// Все переменные имеют префикс UNITRACK_. Ошибки собираются и возвращаются
// одной ошибкой, чтобы администратор увидел сразу все проблемы, а не по одной.
package config

import (
	"errors"
	"fmt"
	"log/slog"
)

// Config — настройки процесса API. Единственное «глобальное» состояние сервера.
type Config struct {
	// HTTPAddr — адрес, на котором слушает HTTP-сервер. UNITRACK_HTTP_ADDR, по умолчанию ":8080".
	HTTPAddr string
	// DatabaseURL — строка подключения к PostgreSQL. UNITRACK_DATABASE_URL, обязательна.
	DatabaseURL string
	// LogLevel — уровень логов. UNITRACK_LOG_LEVEL: debug|info|warn|error, по умолчанию info.
	LogLevel slog.Level
	// Env — окружение. UNITRACK_ENV: development|production, по умолчанию production.
	Env string
}

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

// Load читает настройки через getenv (обычно os.Getenv; в тестах — замыкание над map).
func Load(getenv func(string) string) (Config, error) {
	var problems []error
	cfg := Config{
		HTTPAddr:    withDefault(getenv("UNITRACK_HTTP_ADDR"), ":8080"),
		DatabaseURL: getenv("UNITRACK_DATABASE_URL"),
		LogLevel:    slog.LevelInfo,
		Env:         withDefault(getenv("UNITRACK_ENV"), EnvProduction),
	}

	if cfg.DatabaseURL == "" {
		problems = append(problems, errors.New("UNITRACK_DATABASE_URL: переменная обязательна"))
	}

	if raw := getenv("UNITRACK_LOG_LEVEL"); raw != "" {
		level, err := parseLogLevel(raw)
		if err != nil {
			problems = append(problems, fmt.Errorf("UNITRACK_LOG_LEVEL: %w", err))
		}
		cfg.LogLevel = level
	}

	if cfg.Env != EnvDevelopment && cfg.Env != EnvProduction {
		problems = append(problems, fmt.Errorf("UNITRACK_ENV: допустимы %q или %q, получено %q", EnvDevelopment, EnvProduction, cfg.Env))
	}

	if len(problems) > 0 {
		return Config{}, fmt.Errorf("конфигурация: %w", errors.Join(problems...))
	}
	return cfg, nil
}

func withDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func parseLogLevel(raw string) (slog.Level, error) {
	switch raw {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("допустимы debug|info|warn|error, получено %q", raw)
	}
}
