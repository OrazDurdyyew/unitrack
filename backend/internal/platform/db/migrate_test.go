package db_test

import (
	"context"
	"testing"

	"github.com/OrazDurdyyew/unitrack/backend/internal/platform/db"
	"github.com/OrazDurdyyew/unitrack/backend/internal/platform/testdb"
)

func TestMigrate_AppliesInitialSchema(t *testing.T) {
	pool := testdb.New(t)
	var value string
	err := pool.QueryRow(context.Background(), `select value from app_info where key = 'schema'`).Scan(&value)
	if err != nil {
		t.Fatal(err)
	}
	if value != "unitrack" {
		t.Errorf("value = %q, ожидалось unitrack", value)
	}
}

func TestMigrate_IsIdempotent(t *testing.T) {
	pool := testdb.New(t)
	if err := db.Migrate(context.Background(), pool); err != nil {
		t.Fatalf("повторная миграция: %v", err)
	}
}

func TestNew_IsolatesDatabasesBetweenCalls(t *testing.T) {
	a := testdb.New(t)
	b := testdb.New(t)
	if _, err := a.Exec(context.Background(), `insert into app_info (key, value) values ('only_a', '1')`); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := b.QueryRow(context.Background(), `select count(*) from app_info where key = 'only_a'`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("базы не изолированы: строка из a видна в b")
	}
}
