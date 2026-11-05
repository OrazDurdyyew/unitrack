package config_test

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/OrazDurdyyew/unitrack/backend/internal/platform/config"
)

func envFrom(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoad_Defaults(t *testing.T) {
	cfg, err := config.Load(envFrom(map[string]string{"UNITRACK_DATABASE_URL": "postgres://u:p@h/db"}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, ожидалось :8080", cfg.HTTPAddr)
	}
	if cfg.LogLevel != slog.LevelInfo {
		t.Errorf("LogLevel = %v, ожидалось info", cfg.LogLevel)
	}
	if cfg.Env != "production" {
		t.Errorf("Env = %q, ожидалось production", cfg.Env)
	}
	if cfg.DatabaseURL != "postgres://u:p@h/db" {
		t.Errorf("DatabaseURL = %q", cfg.DatabaseURL)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	_, err := config.Load(envFrom(nil))
	if err == nil || !strings.Contains(err.Error(), "UNITRACK_DATABASE_URL") {
		t.Fatalf("ожидалась ошибка про UNITRACK_DATABASE_URL, получено %v", err)
	}
}

func TestLoad_InvalidLogLevel(t *testing.T) {
	_, err := config.Load(envFrom(map[string]string{
		"UNITRACK_DATABASE_URL": "postgres://u:p@h/db",
		"UNITRACK_LOG_LEVEL":    "loud",
	}))
	if err == nil || !strings.Contains(err.Error(), "UNITRACK_LOG_LEVEL") {
		t.Fatalf("ожидалась ошибка про UNITRACK_LOG_LEVEL, получено %v", err)
	}
}

func TestLoad_InvalidEnv(t *testing.T) {
	_, err := config.Load(envFrom(map[string]string{
		"UNITRACK_DATABASE_URL": "postgres://u:p@h/db",
		"UNITRACK_ENV":          "staging",
	}))
	if err == nil || !strings.Contains(err.Error(), "UNITRACK_ENV") {
		t.Fatalf("ожидалась ошибка про UNITRACK_ENV, получено %v", err)
	}
}

func TestLoad_ReportsAllProblemsAtOnce(t *testing.T) {
	_, err := config.Load(envFrom(map[string]string{"UNITRACK_LOG_LEVEL": "loud"}))
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
	msg := err.Error()
	if !strings.Contains(msg, "UNITRACK_DATABASE_URL") || !strings.Contains(msg, "UNITRACK_LOG_LEVEL") {
		t.Fatalf("ожидались обе проблемы в одной ошибке, получено %q", msg)
	}
}
