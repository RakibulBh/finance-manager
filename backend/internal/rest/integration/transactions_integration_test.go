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

func TestTransactions_Integration(t *testing.T) {
	userRepo := postgres.NewUserRepository(testDB)
	accountRepo := postgres.NewAccountRepository(testDB)
	ledgerRepo := postgres.NewLedgerRepository(testDB)
	authHandler := rest.NewAuthHandler(userRepo, testCfg.JWTSecret)
	accountHandler := rest.NewAccountHandler(accountRepo)
	transactionHandler := rest.NewTransactionHandler(ledgerRepo)

	router := rest.NewRouter(rest.RouterConfig{
		AuthHandler:        authHandler,
		AccountHandler:     accountHandler,
		TransactionHandler: transactionHandler,
		JWTSecret:          testCfg.JWTSecret,
	})

	server := httptest.NewServer(router)
	defer server.Close()

	ClearDB()

	// Setup: User and Account
	DoRequest(server, "POST", "/api/register", `{"email": "tx@example.com", "password": "password123", "family_name": "TX Family"}`, "")
	loginResp, _ := DoRequest(server, "POST", "/api/login", `{"email": "tx@example.com", "password": "password123"}`, "")
	var loginResult map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	token := loginResult["data"].(map[string]interface{})["token"].(string)

	accResp, _ := DoRequest(server, "POST", "/api/accounts", `{"name": "Checking", "balance": 1000, "currency": "USD", "type": "depository", "subtype": "checking"}`, token)
	var account map[string]interface{}
	json.NewDecoder(accResp.Body).Decode(&account)
	accountID := account["id"].(string)

	t.Run("Create Transaction", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{"account_id": "%s", "amount": 50, "date": "%s", "name": "Lunch", "merchant_name": "Cafe"}`, accountID, time.Now().Format(time.RFC3339))
		resp, err := DoRequest(server, "POST", "/api/transactions", reqBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var tx map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&tx)
		assert.Equal(t, "Lunch", tx["name"])
		assert.Equal(t, float64(50), tx["amount"])
	})

	t.Run("Create Transfer", func(t *testing.T) {
		// Need another account
		accResp2, _ := DoRequest(server, "POST", "/api/accounts", `{"name": "Savings", "balance": 0, "currency": "USD", "type": "depository", "subtype": "savings"}`, token)
		var account2 map[string]interface{}
		json.NewDecoder(accResp2.Body).Decode(&account2)
		accountID2 := account2["id"].(string)

		reqBody := fmt.Sprintf(`{"from_account_id": "%s", "to_account_id": "%s", "amount": 200, "date": "%s", "name": "Savings Transfer"}`, accountID, accountID2, time.Now().Format(time.RFC3339))
		resp, err := DoRequest(server, "POST", "/api/transfers", reqBody, token)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}
