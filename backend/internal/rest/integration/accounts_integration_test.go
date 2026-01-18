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

func TestAccounts_Integration(t *testing.T) {
	// Initialize handlers
	userRepo := postgres.NewUserRepository(testDB)
	accountRepo := postgres.NewAccountRepository(testDB)
	authHandler := rest.NewAuthHandler(userRepo, testCfg.JWTSecret)
	accountHandler := rest.NewAccountHandler(accountRepo)

	router := rest.NewRouter(rest.RouterConfig{
		AuthHandler:    authHandler,
		AccountHandler: accountHandler,
		JWTSecret:      testCfg.JWTSecret,
	})

	server := httptest.NewServer(router)
	defer server.Close()

	ClearDB()

	// 1. Register and Login to get token
	DoRequest(server, "POST", "/api/register", `{"email": "account@example.com", "password": "password123", "family_name": "Account Family"}`, "")
	resp, _ := DoRequest(server, "POST", "/api/login", `{"email": "account@example.com", "password": "password123"}`, "")
	var loginResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&loginResult)
	token := loginResult["data"].(map[string]interface{})["token"].(string)

	t.Run("Create Account Success", func(t *testing.T) {
		reqBody := `{"name": "Savings", "balance": 1000, "currency": "USD", "type": "depository", "subtype": "savings"}`
		resp, err := DoRequest(server, "POST", "/api/accounts", reqBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var acc map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&acc)
		assert.Equal(t, "Savings", acc["name"])
		assert.Equal(t, "asset", acc["classification"])
	})

	t.Run("List Accounts", func(t *testing.T) {
		resp, err := DoRequest(server, "GET", "/api/accounts", "", token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		data := result["data"].(map[string]interface{})
		accounts := data["accounts"].([]interface{})
		assert.True(t, len(accounts) >= 1)
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		resp, err := DoRequest(server, "GET", "/api/accounts", "", "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
