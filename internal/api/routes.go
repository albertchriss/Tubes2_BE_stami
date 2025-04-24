package api

import (
	_ "github.com/albertchriss/Tubes2_BE_stami/docs"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the API routes for the non-versioned endpoints.
func RegisterRoutes(r *gin.Engine, handlers *Handlers) {
	r.GET("/docs/*any", handlers.DocsHandler)
	r.GET("/health", handlers.HealthHandler.HealthCheck)
}
