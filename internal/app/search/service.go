package search

import (
	"log"

	"github.com/albertchriss/Tubes2_BE_stami/internal/app/socket"
	"github.com/albertchriss/Tubes2_BE_stami/internal/core"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
	"github.com/albertchriss/Tubes2_BE_stami/internal/utils"
)

type Service interface {
	BFSSearch(query string, numRecipe int, liveUpdate bool) scraper.TreeNode
	DFSSearch(query string, numRecipe int, liveUpdate bool) scraper.TreeNode
	BidirectionalSearch(query string, numMeetingNodeChoice int) scraper.TreeNode
}

type service struct {
	appCtx    *core.AppContext
	wsManager *socket.ClientManager
}

func NewService(appCtx *core.AppContext, wsManager *socket.ClientManager) *service {
	return &service{
		appCtx:    appCtx,
		wsManager: wsManager,
	}
}

func (s *service) BFSSearch(query string, numRecipe int, liveUpdate bool) scraper.TreeNode {
	log.Println("Performing BFS search for query:", query)
	recipe := s.appCtx.Config.RecipeTree
	if recipe == nil {
		log.Println("Recipe tree is nil")
		return scraper.TreeNode{Name: "Recipe tree is nil"}
	}
	if _, exists := (*recipe)[query]; !exists {
		log.Println("Query not found in recipe tree")
		return scraper.TreeNode{Name: "Query not found in recipe tree"}
	}

	if numRecipe > 1 {
		log.Println("Performing BFS for multiple recipes")
		return utils.MultipleRecipeBFS(recipe, s.appCtx.Config.TierMap, query, numRecipe, liveUpdate, s.wsManager)
	} else {
		log.Println("Performing BFS for single recipe")
		return utils.SingleRecipeBFS(recipe, s.appCtx.Config.TierMap, query, liveUpdate, s.wsManager)
	}
}

func (s *service) DFSSearch(query string, numRecipe int, liveUpdate bool) scraper.TreeNode {
	log.Println("Performing DFS search for query:", query)
	recipe := s.appCtx.Config.RecipeTree
	if recipe == nil {
		log.Println("Recipe tree is nil")
		return scraper.TreeNode{Name: "Recipe tree is nil"}
	}
	if _, exists := (*recipe)[query]; !exists {
		log.Println("Query not found in recipe tree")
		return scraper.TreeNode{Name: "Query not found in recipe tree"}
	}

	if numRecipe > 1 {
		log.Println("Performing DFS for multiple recipes")
		return utils.MultipleRecipeDFS(recipe, s.appCtx.Config.TierMap, query, numRecipe, liveUpdate, s.wsManager)
	} else {
		log.Println("Performing DFS for single recipe")
		return utils.SingleRecipeDFS(recipe, s.appCtx.Config.TierMap, query, liveUpdate, s.wsManager)
	}
}

func (s *service) BidirectionalSearch(query string, numMeetingNodeChoice int) scraper.TreeNode {
	log.Println("Performing Bidirectional search for query:", query, "meeting node choice:", numMeetingNodeChoice)
	recipe := s.appCtx.Config.RecipeTree
	tierMap := s.appCtx.Config.TierMap

	if recipe == nil {
		log.Println("Recipe tree is nil")
		return scraper.TreeNode{Name: "Recipe tree is nil"}
	}

	return utils.BidirectionalSearch(recipe, tierMap, query, numMeetingNodeChoice)
}
