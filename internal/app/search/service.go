package search

import (
	"log"

	"github.com/albertchriss/Tubes2_BE_stami/internal/core"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
	"github.com/albertchriss/Tubes2_BE_stami/internal/utils"
)

type Service interface {
	BFSSearch(query string, tipe string) scraper.TreeNode
}

type service struct {
	appCtx *core.AppContext
}

func NewService(appCtx *core.AppContext) *service {
	return &service{
		appCtx: appCtx,
	}
}

func (s *service) BFSSearch(query string, tipe string) scraper.TreeNode {
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

	// if (tipe == scraper.SINGLERECIPE)
	if tipe == scraper.MULTIPLERECIPE {
		log.Println("Performing BFS for multiple recipes")
		return utils.MultipleRecipeBFS(*recipe, query)
	}

	return scraper.TreeNode{Name: "Invalid recipe type"}

}
