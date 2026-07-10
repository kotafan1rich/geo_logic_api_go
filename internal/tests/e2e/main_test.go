//go:build e2e

package e2e

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kotafan1rich/geo_logic_api_go/internal/app"
	"github.com/kotafan1rich/geo_logic_api_go/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testDBName     = "geo_logic_test"
	testDBUser     = "test_user"
	testDBPassword = "test_pass"
)

var (
	testServerURL string
	testDB        *sql.DB
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgis/postgis:18-3.6",
		postgres.WithDatabase(testDBName),
		postgres.WithUsername(testDBUser),
		postgres.WithPassword(testDBPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second),
		),
	)

	if err != nil {
		log.Fatalf("Не удалось запустить postgres в docker: %s", err)
	}

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", testDBUser)
	os.Setenv("DB_PASSWORD", testDBPassword)
	os.Setenv("DB_NAME", testDBName)

	os.Setenv("SERVER_PORT", "8081")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", testDBUser, testDBPassword, host, port.Port(), testDBName)

	testDB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к тестовой БД: %s", err)
	}

	go func() {
		config.MustLoad()
		app := app.New()
		if err := app.Run(); err != nil {
			slog.Error("app error", "err", err)
			os.Exit(1)
		}
	}()

	testServerURL = "http://localhost:8081"
	time.Sleep(500 * time.Millisecond)

	code := m.Run()

	testDB.Close()

	_ = pgContainer.Terminate(ctx)

	os.Exit(code)
}

func usersAPI() string { return testServerURL + "/api/users" }

func rentsAPI() string { return testServerURL + "/api/rents" }

func parseBody(body io.ReadCloser, dest any) error {
	err := json.NewDecoder(body).Decode(dest)
	if err != nil {
		return err
	}
	return nil
}

func clearTables(t *testing.T) {
	_, err := testDB.Exec("TRUNCATE TABLE users, rents RESTART IDENTITY CASCADE;")
	if err != nil {
		t.Fatalf("Ошибка очистки таблицы перед тестом: %v", err)
	}
}
