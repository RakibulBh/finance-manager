package integration

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	testDB    *pgxpool.Pool
	testRedis *redis.Client
	testCfg   *config.Config
)

// SetupTestEnv initializes the test database and redis
func SetupTestEnv() {
	log.Println("Setting up test environment...")
	// Load config
	testCfg = &config.Config{
		DatabaseURL: os.Getenv("TEST_DATABASE_URL"),
		RedisAddr:   os.Getenv("TEST_REDIS_ADDR"),
		JWTSecret:   "test-secret",
	}

	if testCfg.DatabaseURL == "" {
		testCfg.DatabaseURL = "postgres://user:password@localhost:5435/finance_manager_test?sslmode=disable"
	}
	if testCfg.RedisAddr == "" {
		testCfg.RedisAddr = "localhost:6379"
	}

	// Run migrations
	runMigrations(testCfg.DatabaseURL)

	// Connect to DB
	var err error
	testDB, err = pgxpool.New(context.Background(), testCfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}

	// Connect to Redis
	testRedis = redis.NewClient(&redis.Options{
		Addr: testCfg.RedisAddr,
	})
	if err := testRedis.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to test Redis: %v", err)
	}
}

func runMigrations(dbURL string) {
	// Find migrations directory
	cwd, _ := os.Getwd()
	// Navigate up to find db/migrations from internal/rest/integration
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "db/migrations")); err == nil {
			dir = filepath.Join(dir, "db/migrations")
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal("Could not find db/migrations directory")
		}
		dir = parent
	}

	m, err := migrate.New("file://"+dir, dbURL)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration Up failed: %v", err)
	}
}

// TeardownTestEnv cleans up the test environment
func TeardownTestEnv() {
	if testDB != nil {
		testDB.Close()
	}
	if testRedis != nil {
		testRedis.FlushAll(context.Background())
		testRedis.Close()
	}
}

// ClearDB removes all data from the test database
func ClearDB() {
	tables := []string{"trades", "securities", "entries", "transactions", "accounts", "users", "families"}
	for _, table := range tables {
		_, err := testDB.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			log.Printf("Failed to truncate %s: %v", table, err)
		}
	}
}

// DoRequest performs a request against a test server
func DoRequest(server *httptest.Server, method, path string, body string, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, server.URL+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return http.DefaultClient.Do(req)
}

// TestMain runs before all tests in the package
func TestMain(m *testing.M) {
	SetupTestEnv()
	code := m.Run()
	TeardownTestEnv()
	os.Exit(code)
}
