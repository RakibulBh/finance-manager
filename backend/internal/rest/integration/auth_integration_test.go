package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rakibulbh/ai-finance-manager/internal/repository/postgres"
	"github.com/rakibulbh/ai-finance-manager/internal/rest"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Integration(t *testing.T) {
	// Initialize handlers with real repositories
	userRepo := postgres.NewUserRepository(testDB)
	authHandler := rest.NewAuthHandler(userRepo, testCfg.JWTSecret)

	router := rest.NewRouter(rest.RouterConfig{
		AuthHandler: authHandler,
		JWTSecret:   testCfg.JWTSecret,
	})

	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("Register Success", func(t *testing.T) {
		ClearDB()
		reqBody := `{"email": "test@example.com", "password": "SecurePass123!", "family_name": "Test Family"}`
		resp, err := DoRequest(server, "POST", "/api/register", reqBody, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result[0]["data"].(map[string]interface{})
		assert.Equal(t, "test@example.com", data["email"])
	})

	t.Run("Login Success", func(t *testing.T) {
		// User already registered from previous run if we didn't clear
		// But we cleared in the subtest. Let's register again.
		ClearDB()
		DoRequest(server, "POST", "/api/register", `{"email": "login@example.com", "password": "password123"}`, "")

		reqBody := `{"email": "login@example.com", "password": "password123"}`
		resp, err := DoRequest(server, "POST", "/api/login", reqBody, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		assert.NotNil(t, data["token"])
		assert.Equal(t, "login@example.com", data["user"].(map[string]interface{})["email"])
	})

	t.Run("Login Failure - Invalid Credentials", func(t *testing.T) {
		reqBody := `{"email": "login@example.com", "password": "wrongpassword"}`
		resp, err := DoRequest(server, "POST", "/api/login", reqBody, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
