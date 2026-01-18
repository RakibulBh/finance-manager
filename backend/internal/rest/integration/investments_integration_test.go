package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rakibulbh/ai-finance-manager/internal/repository/postgres"
	"github.com/rakibulbh/ai-finance-manager/internal/rest"
	"github.com/stretchr/testify/assert"
)

func TestInvestments_Integration(t *testing.T) {
	userRepo := postgres.NewUserRepository(testDB)
	accountRepo := postgres.NewAccountRepository(testDB)
	investmentRepo := postgres.NewInvestmentRepository(testDB)
	authHandler := rest.NewAuthHandler(userRepo, testCfg.JWTSecret)
	accountHandler := rest.NewAccountHandler(accountRepo)
	investmentHandler := rest.NewInvestmentHandler(investmentRepo)

	router := rest.NewRouter(rest.RouterConfig{
		AuthHandler:       authHandler,
		AccountHandler:    accountHandler,
		InvestmentHandler: investmentHandler,
		JWTSecret:         testCfg.JWTSecret,
	})

	server := httptest.NewServer(router)
	defer server.Close()

	ClearDB()

	// Setup: User and Account
	DoRequest(server, "POST", "/api/register", `{"email": "invest@example.com", "password": "password123", "family_name": "Invest Family"}`, "")
	loginResp, _ := DoRequest(server, "POST", "/api/login", `{"email": "invest@example.com", "password": "password123"}`, "")
	var loginResult map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	token := loginResult["data"].(map[string]interface{})["token"].(string)

	accResp, _ := DoRequest(server, "POST", "/api/accounts", `{"name": "Brokerage", "balance": 5000, "currency": "USD", "type": "investment", "subtype": "brokerage"}`, token)
	var account map[string]interface{}
	json.NewDecoder(accResp.Body).Decode(&account)
	accountID := account["id"].(string)

	t.Run("Create Trade", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{
			"account_id": "%s",
			"ticker": "AAPL",
			"security_name": "Apple Inc.",
			"qty": 10,
			"price": 150.50,
			"date": "%s",
			"kind": "buy"
		}`, accountID, time.Now().Format(time.RFC3339))

		resp, err := DoRequest(server, "POST", "/api/investments/trade", reqBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "buy AAPL", result["name"])
		assert.Equal(t, -1505.0, result["amount"])
	})
}
