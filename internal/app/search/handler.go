package search

import (
	"net/http"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
	"github.com/gin-gonic/gin"
)

type SearchResponse struct {
	Message string           `json:"message"`
	Result  scraper.TreeNode `json:"result"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// BFSSearchHandler godoc
// @Summary BFS search handler
// @Description Search the recipe of elements using BFS
// @Tags Search
// @Accept json
// @Produce json
// @Param q query string true "Query parameter"
// @Param tipe query string true "Search type" enums(single, multiple)
// @Success 200 {object} SearchResponse
// @Router /search/bfs [get]
func (h *Handler) BFSSearchHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "Query parameter is required",
		})
		return
	}

	tipe := c.Query("tipe")
	if tipe == "" {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "Type parameter is required",
		})
		return
	}

	if tipe != scraper.SINGLERECIPE && tipe != scraper.MULTIPLERECIPE {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "Type parameter must be single or multiple",
		})
		return
	}

	res := h.service.DFSSearch(query, tipe)
	c.JSON(http.StatusOK, SearchResponse{
		Message: "BFS search completed",
		Result:  res,
	})
}

// DFSSearchHandler godoc
// @Summary DFS search handler
// @Description Search the recipe of elements using DFS
// @Tags Search
// @Accept json
// @Produce json
// @Param q query string true "Query parameter"
// @Param tipe query string true "Search type" enums(single, multiple)
// @Success 200 {object} SearchResponse
// @Router /search/dfs [get]
func (h *Handler) DFSSearchHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "Query parameter is required",
		})
		return
	}

	tipe := c.Query("tipe")
	if tipe == "" {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "Type parameter is required",
		})
		return
	}

	if tipe != scraper.SINGLERECIPE && tipe != scraper.MULTIPLERECIPE {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "Type parameter must be single or multiple",
		})
		return
	}

	res := h.service.BFSSearch(query, tipe)
	c.JSON(http.StatusOK, SearchResponse{
		Message: "BFS search completed",
		Result:  res,
	})
}
