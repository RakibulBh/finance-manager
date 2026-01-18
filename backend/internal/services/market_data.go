package services

import (
	"math/rand"
)

type MockMarketData struct{}

func (m *MockMarketData) GetQuote(ticker string) (float64, error) {
	// Return a random price between 10 and 1000 for demonstration
	return 10 + rand.Float64()*(1000-10), nil
}

func NewMockMarketData() *MockMarketData {
	return &MockMarketData{}
}
