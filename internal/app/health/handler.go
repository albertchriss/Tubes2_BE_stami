package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Check if the application is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:  "success",
		Message: "Application is running and healthy",
	})
}
