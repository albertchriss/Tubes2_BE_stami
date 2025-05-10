package search

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/albertchriss/Tubes2_BE_stami/internal/core"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
	"github.com/gin-gonic/gin"
)

type ElementResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type SearchResponse struct {
	Message string           `json:"message"`
	Result  scraper.TreeNode `json:"result"`
}

type Handler struct {
	service Service
	AppCtx  *core.AppContext
}

func NewHandler(service Service, appCtx *core.AppContext) *Handler {
	return &Handler{
		service: service,
		AppCtx:  appCtx,
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

// BidirectionalSearchHandler godoc
// @Summary Bidirectional search handler
// @Description Search the recipe of elements using Bidirectional Search.
// @Tags Search
// @Accept json
// @Produce json
// @Param q query string true "Target element to search for"
// @Param num query string false "Chooses the Nth found meeting node (sorted) to construct the path" default(1)
// @Success 200 {object} SearchResponse "Successful search operation."
// @Router /search/bidirectional [get]
func (h *Handler) BidirectionalSearchHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "Query parameter is required",
		})
		return
	}

	numChoiceStr := c.DefaultQuery("num", "1")
	numChoiceInt, err := strconv.Atoi(numChoiceStr)
	if err != nil || numChoiceInt < 1 {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Message: "num parameter must be a positive integer",
		})
		return
	}

	res := h.service.BidirectionalSearch(query, numChoiceInt)
	c.JSON(http.StatusOK, SearchResponse{
		Message: "Bidirectional search completed",
		Result:  res,
	})
}

// GetElementsHandler godoc
// @Summary Get all elements
// @Description Get a list of all available elements
// @Tags Elements
// @Produce json
// @Success 200 {array} ElementResponse
// @Router /elements [get]
func (h *Handler) GetElementsHandler(c *gin.Context) {
	tierData := h.AppCtx.Config.TierMap
	if tierData == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tier data not loaded"})
		return
	}

	var elements []ElementResponse
	for elementName := range *tierData {
		elements = append(elements, ElementResponse{Value: elementName, Label: elementName})
	}

	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Label < elements[j].Label
	})

	c.JSON(http.StatusOK, elements)
}