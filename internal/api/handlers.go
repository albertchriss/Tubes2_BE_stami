package api

import (
	"github.com/albertchriss/Tubes2_BE_stami/internal/app/health"
	"github.com/albertchriss/Tubes2_BE_stami/internal/app/search"
	"github.com/albertchriss/Tubes2_BE_stami/internal/core"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handlers is a struct that contains all the handlers for the application.
type Handlers struct {
	DocsHandler   gin.HandlerFunc
	HealthHandler *health.Handler
	SearchHandler *search.Handler
}

// InitHandlers initializes all the handlers for the application.
// It takes an AppContext as a parameter and returns a Handlers struct.
// The AppContext contains all the app dependencies such as the database connection pool and Redis client.
func InitHandlers(appCtx *core.AppContext) *Handlers {
	// Docs Handler Initialization
	docsHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	// Health Handler Initialization
	healthHandler := health.NewHandler()

	searchService := search.NewService(appCtx)
	searchHandler := search.NewHandler(searchService)

	return &Handlers{
		DocsHandler:   docsHandler,
		HealthHandler: healthHandler,
		SearchHandler: searchHandler,
	}
}
