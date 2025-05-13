package app

import (
	"github.com/albertchriss/Tubes2_BE_stami/docs"
	"github.com/albertchriss/Tubes2_BE_stami/internal/api"
	"github.com/albertchriss/Tubes2_BE_stami/internal/core"
	"github.com/albertchriss/Tubes2_BE_stami/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run initializes the application context and starts the HTTP server.
// It sets up all necessary components, such as the database connection pool.
// It also registers the API routes.
func Run() {
	r := gin.Default()

	cfg := core.NewAppConfig()

	corsConfig := cors.DefaultConfig()
	allowedOrigin := utils.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:3000")
	corsConfig.AllowOrigins = []string{allowedOrigin}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	// on production
	docs.SwaggerInfo.BasePath = "/api"

	appCtx := core.AppContext{
		Config: cfg,
	}

	handlers := api.InitHandlers(&appCtx)

	r.Use(cors.New(corsConfig))
	
	api.RegisterRoutes(r, handlers, &appCtx)

	r.Run(cfg.AppAddress)
}
