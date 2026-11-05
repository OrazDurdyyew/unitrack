// Package health — служебные операции: живость сервера и готовность зависимостей.
package health

import (
	"context"
	"net/http"
	"time"

	"github.com/OrazDurdyyew/unitrack/backend/internal/api/gen"
	"github.com/OrazDurdyyew/unitrack/backend/internal/platform/httpserver"
)

// Pinger проверяет доступность зависимости (например, *pgxpool.Pool).
type Pinger interface {
	Ping(ctx context.Context) error
}

// Handler реализует операции health и ready из контракта.
type Handler struct {
	version string
	db      Pinger
}

// Проверка на этапе компиляции: Handler покрывает весь контракт.
var _ gen.ServerInterface = (*Handler)(nil)

const pingTimeout = 2 * time.Second

// New создаёт обработчик. version — версия сборки, db — проверка базы данных.
func New(version string, db Pinger) *Handler {
	return &Handler{version: version, db: db}
}

// GetHealth — сервер запущен. Базу не проверяет.
func (h *Handler) GetHealth(w http.ResponseWriter, _ *http.Request) {
	httpserver.WriteJSON(w, http.StatusOK, gen.HealthResponse{
		Status:  gen.HealthResponseStatusOk,
		Version: h.version,
	})
}

// GetReady — готовность обслуживать запросы: база должна отвечать.
func (h *Handler) GetReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), pingTimeout)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		httpserver.WriteJSON(w, http.StatusServiceUnavailable, gen.ReadyResponse{
			Status:   gen.ReadyResponseStatusDegraded,
			Database: gen.ReadyResponseDatabaseUnavailable,
		})
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, gen.ReadyResponse{
		Status:   gen.ReadyResponseStatusOk,
		Database: gen.ReadyResponseDatabaseOk,
	})
}
