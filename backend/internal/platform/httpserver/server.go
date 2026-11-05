// Package httpserver — маршрутизация, промежуточные обработчики и формат ошибок HTTP.
package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/OrazDurdyyew/unitrack/backend/internal/api/gen"
)

// BasePath — префикс всех операций контракта (servers[0].url в api/openapi.yaml).
const BasePath = "/api/v1"

// NewRouter собирает маршрутизатор: промежуточные обработчики, операции из контракта,
// ошибки 404/405 в формате problem+json.
func NewRouter(log *slog.Logger, h gen.ServerInterface) http.Handler {
	r := chi.NewRouter()
	r.Use(RequestID)
	r.Use(Recover(log))
	r.Use(Logger(log))

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		WriteProblem(w, req, http.StatusNotFound, "not_found", "Ресурс не найден", "")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		WriteProblem(w, req, http.StatusMethodNotAllowed, "method_not_allowed", "Метод не поддерживается", "")
	})

	gen.HandlerWithOptions(h, gen.ChiServerOptions{
		BaseURL:    BasePath,
		BaseRouter: r,
		ErrorHandlerFunc: func(w http.ResponseWriter, req *http.Request, err error) {
			WriteProblem(w, req, http.StatusBadRequest, "bad_request", "Некорректный запрос", err.Error())
		},
	})
	return r
}
