package app

import (
	"github.com/albertchriss/Tubes2_BE_stami/internal/api"
	"github.com/albertchriss/Tubes2_BE_stami/internal/core"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	config := cors.DefaultConfig()
	// Izinkan origin frontend Anda. Untuk pengembangan, bisa localhost:3000
	config.AllowOrigins = []string{"http://localhost:3000"}
	// Metode HTTP yang diizinkan
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// Header HTTP yang diizinkan
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(config))

	api.RegisterRoutes(r, handlers, &appCtx)

	r.Run(cfg.AppAddress)
}
