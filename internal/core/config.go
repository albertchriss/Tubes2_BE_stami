package core

import (
	"log"

	"github.com/albertchriss/Tubes2_BE_stami/internal/utils"
	"github.com/joho/godotenv"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

// AppConfig holds all the application configuration
type AppConfig struct {
	AppName    string
	AppAddress string
	AppPort    string
	RecipeTree *scraper.Recipe
}

// NewAppConfig initializes the application configuration
func NewAppConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ WARNING: Could not load .env file")
	}

	outputFile := "data/little_alchemy_recipes.json"
	scraper.Scraper(outputFile)

	appName := utils.GetString("APP_NAME", "")
	appAddress := utils.GetString("APP_ADDRESS", "")
	appPort := utils.GetString("APP_PORT", "8080")
	recipeTree := scraper.JsonToMap(outputFile)

	return &AppConfig{
		AppName:    appName,
		AppAddress: appAddress,
		AppPort:    appPort,
		RecipeTree: recipeTree,
	}
}
