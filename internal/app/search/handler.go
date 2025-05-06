package search

import (
	"net/http"
	"strconv"

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
// @Param num query string false "Number of recipes to return" default(1)
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

	numRecipe := c.Query("num")
	if numRecipe == "" {
		numRecipe = "1"
	}
	numRecipeInt, err := strconv.Atoi(numRecipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "num parameter must be an integer",
		})
		return
	}

	if numRecipeInt < 1 {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "num parameter must be greater than 0",
		})
		return
	}

	res := h.service.BFSSearch(query, numRecipeInt)
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
// @Param num query string false "Number of recipes to return" default(1)
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

	numRecipe := c.Query("num")
	if numRecipe == "" {
		numRecipe = "1"
	}
	numRecipeInt, err := strconv.Atoi(numRecipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "num parameter must be an integer",
		})
		return
	}

	if numRecipeInt < 1 {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "num parameter must be greater than 0",
		})
		return
	}

	res := h.service.DFSSearch(query, numRecipeInt)
	c.JSON(http.StatusOK, SearchResponse{
		Message: "DFS search completed",
		Result:  res,
	})
}
