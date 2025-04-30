package search

import (
	"log"

	"github.com/albertchriss/Tubes2_BE_stami/internal/core"
)

type Service interface {
	BFSSearch(query string) []string
}

type service struct {
	appCtx *core.AppContext
}

func NewService(appCtx *core.AppContext) *service {
	return &service{
		appCtx: appCtx,
	}
}

func (s *service) BFSSearch(query string) []string {
	log.Println("Performing BFS search for query:", query)
	recipe := s.appCtx.Config.RecipeTree
	if recipe == nil {
		log.Println("Recipe tree is nil")
		return []string{"Recipe tree is nil"}
	}
	if _, exists := (*recipe)[query]; !exists {
		log.Println("Query not found in recipe tree")
		return []string{"Query not found in recipe tree"}
	}
	var str []string
	arr := (*recipe)[query]
	for i := range arr {
		str = append(str, arr[i].First()+" - "+arr[i].Second())
	}

	// visited := make(map[string]bool)
	// queue := []string{query}
	// visited[query] = true
	return str

}
