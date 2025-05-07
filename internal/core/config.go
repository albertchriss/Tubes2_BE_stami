package core

import (
	"log"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
	"github.com/albertchriss/Tubes2_BE_stami/internal/utils"
	"github.com/joho/godotenv"
)

// AppConfig holds all the application configuration
type AppConfig struct {
	AppName    string
	AppAddress string
	AppPort    string
	RecipeTree *scraper.Recipe
	TierMap    *scraper.Tier
}

// NewAppConfig initializes the application configuration
func NewAppConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ WARNING: Could not load .env file")
	}

	recipeFile := "data/little_alchemy_recipes.json"
	tierFile := "data/little_alchemy_tiers.json"
	scraper.Scraper(recipeFile, tierFile)

	appName := utils.GetString("APP_NAME", "")
	appAddress := utils.GetString("APP_ADDRESS", "")
	appPort := utils.GetString("APP_PORT", "8080")
	recipeTree := scraper.JsonToRecipe(recipeFile)
	tierMap := scraper.JsonToTier(tierFile)

	recipeTree.SortRecipeChildren(tierMap)

	return &AppConfig{
		AppName:    appName,
		AppAddress: appAddress,
		AppPort:    appPort,
		RecipeTree: recipeTree,
		TierMap:    tierMap,
	}
}
