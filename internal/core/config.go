package core

import (
	"log"

	"github.com/albertchriss/Tubes2_BE_stami/internal/utils"
	"github.com/joho/godotenv"
)

// AppConfig holds all the application configuration
type AppConfig struct {
	AppName    string
	AppAddress string
	AppPort    string
}

// NewAppConfig initializes the application configuration
func NewAppConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ WARNING: Could not load .env file")
	}

	appName := utils.GetString("APP_NAME", "")
	appAddress := utils.GetString("APP_ADDRESS", "")
	appPort := utils.GetString("APP_PORT", "8080")

	return &AppConfig{
		AppName:    appName,
		AppAddress: appAddress,
		AppPort:    appPort,
	}
}
