package health_test

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OrazDurdyyew/unitrack/backend/internal/api/gen"
	"github.com/OrazDurdyyew/unitrack/backend/internal/health"
	"github.com/OrazDurdyyew/unitrack/backend/internal/platform/httpserver"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

func newServer(t *testing.T, pinger health.Pinger) *httptest.Server {
	t.Helper()
	h := health.New("1.2.3", pinger)
	srv := httptest.NewServer(httpserver.NewRouter(slog.New(slog.DiscardHandler), h))
	t.Cleanup(srv.Close)
	return srv
}

func get(t *testing.T, url string) *http.Response {
	t.Helper()
	res, err := http.Get(url) //nolint:gosec // адрес тестового сервера
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = res.Body.Close() })
	return res
}

func TestHealth_ReturnsOkAndVersion(t *testing.T) {
	srv := newServer(t, fakePinger{})
	res := get(t, srv.URL+"/api/v1/health")
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var body gen.HealthResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Status != gen.HealthResponseStatusOk || body.Version != "1.2.3" {
		t.Errorf("body = %+v", body)
	}
}

func TestReady_DatabaseUp(t *testing.T) {
	srv := newServer(t, fakePinger{})
	res := get(t, srv.URL+"/api/v1/ready")
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var body gen.ReadyResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Status != gen.ReadyResponseStatusOk || body.Database != gen.ReadyResponseDatabaseOk {
		t.Errorf("body = %+v", body)
	}
}

func TestReady_DatabaseDown(t *testing.T) {
	srv := newServer(t, fakePinger{err: errors.New("down")})
	res := get(t, srv.URL+"/api/v1/ready")
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, ожидалось 503", res.StatusCode)
	}
	var body gen.ReadyResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Status != gen.ReadyResponseStatusDegraded || body.Database != gen.ReadyResponseDatabaseUnavailable {
		t.Errorf("body = %+v", body)
	}
}

func TestUnknownRoute_ReturnsProblemJSON(t *testing.T) {
	srv := newServer(t, fakePinger{})
	res := get(t, srv.URL+"/api/v1/nope")
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/problem+json" {
		t.Errorf("Content-Type = %q", ct)
	}
	var body gen.Problem
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Code != "not_found" || body.Status != http.StatusNotFound {
		t.Errorf("body = %+v", body)
	}
}

func TestEveryResponse_HasRequestID(t *testing.T) {
	srv := newServer(t, fakePinger{})
	res := get(t, srv.URL+"/api/v1/health")
	if res.Header.Get("X-Request-Id") == "" {
		t.Error("ожидался заголовок X-Request-Id")
	}
}
