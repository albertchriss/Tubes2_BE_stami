package api

import (
	_ "github.com/albertchriss/Tubes2_BE_stami/docs"
	"github.com/albertchriss/Tubes2_BE_stami/internal/core"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the API routes for the non-versioned endpoints.
func RegisterRoutes(r *gin.Engine, handlers *Handlers, appCtx *core.AppContext) {
	r.GET("/docs/*any", handlers.DocsHandler)
	r.GET("/health", handlers.HealthHandler.HealthCheck)
	searchGroup := r.Group("/search")
	{
		searchGroup.GET("/bfs", handlers.SearchHandler.BFSSearchHandler)
		searchGroup.GET("/dfs", handlers.SearchHandler.DFSSearchHandler)
		searchGroup.GET("/bidirectional", handlers.SearchHandler.BidirectionalSearchHandler)
	}
	r.GET("/ws", handlers.SocketHandler.WebSocketConnectHandler)
	r.GET("/elements", handlers.SearchHandler.GetElementsHandler)
}
