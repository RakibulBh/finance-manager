package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/plaid/plaid-go/v20/plaid"
)

type PlaidService struct {
	client        *plaid.APIClient
	encryptionKey []byte
}

func NewPlaidService(clientID, secret, env, encKey string) *PlaidService {
	conf := plaid.NewConfiguration()
	conf.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	conf.AddDefaultHeader("PLAID-SECRET", secret)

	switch env {
	case "sandbox":
		conf.UseEnvironment(plaid.Sandbox)
	case "development":
		conf.UseEnvironment(plaid.Development)
	case "production":
		conf.UseEnvironment(plaid.Production)
	default:
		conf.UseEnvironment(plaid.Sandbox)
	}

	return &PlaidService{
		client:        plaid.NewAPIClient(conf),
		encryptionKey: []byte(encKey),
	}
}

func (s *PlaidService) CreateLinkToken(ctx context.Context, userID, clientName string) (string, error) {
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: userID,
	}
	request := plaid.NewLinkTokenCreateRequest(
		clientName,
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		user,
	)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_TRANSACTIONS})

	resp, _, err := s.client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}

	return resp.GetLinkToken(), nil
}

func (s *PlaidService) ExchangePublicToken(ctx context.Context, publicToken string) (string, string, error) {
	exchangeRequest := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	resp, _, err := s.client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(*exchangeRequest).Execute()
	if err != nil {
		return "", "", err
	}

	return resp.GetAccessToken(), resp.GetItemId(), nil
}

// Encryption Helpers
func (s *PlaidService) EncryptToken(token string) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *PlaidService) DecryptToken(encryptedToken string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
