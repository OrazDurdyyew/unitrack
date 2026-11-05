package httpserver

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Problem — тело ошибки по RFC 9457. Совпадает со схемой Problem в контракте.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Code     string `json:"code"`
}

// WriteProblem пишет ошибку в формате application/problem+json.
// code — машинный код в snake_case для перевода на клиенте; title/detail — для людей.
func WriteProblem(w http.ResponseWriter, r *http.Request, status int, code, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	p := Problem{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: r.URL.Path,
		Code:     code,
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		slog.Default().Error("не удалось записать ответ об ошибке", "err", err)
	}
}

// WriteJSON пишет успешный ответ в JSON.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Default().Error("не удалось записать ответ", "err", err)
	}
}
