package main

import (
	"fmt"
	"log"

	"github.com/rakibulbh/ai-finance-manager/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}

	fmt.Printf("Starting application in %s mode...\n", cfg.Environment)
	fmt.Printf("Server will run on port %s\n", cfg.Port)
	fmt.Printf("Database URL: %s\n", cfg.DatabaseURL)
}
