package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchResponse struct {
	Message string `json:"message"`
	Result  []string `json:"result"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) BFSSearchHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "Query parameter is required",
		})
		return
	}

	str := h.service.BFSSearch(query)
	c.JSON(http.StatusOK, SearchResponse{
		Message: "BFS search completed",
		Result:  str,
	})
}
