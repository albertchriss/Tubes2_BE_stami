package app

import (
	"github.com/albertchriss/Tubes2_BE_stami/internal/api"
	v1 "github.com/albertchriss/Tubes2_BE_stami/internal/api/v1"
	"github.com/albertchriss/Tubes2_BE_stami/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

// Run initializes the application context and starts the HTTP server.
// It sets up all necessary components, such as the database connection pool.
// It also registers the API routes.
func Run() {
	cfg := core.NewAppConfig()

	appCtx := core.AppContext{
		Config: cfg,
	}

	handlers := api.InitHandlers(&appCtx)

	outputFile := "data/little_alchemy_recipes.json"
	scraper.Scraper(outputFile)

	r := gin.Default()

	api.RegisterRoutes(r, handlers)
	v1.RegisterRoutes(r, handlers, &appCtx)

	r.Run(cfg.AppAddress)
}
