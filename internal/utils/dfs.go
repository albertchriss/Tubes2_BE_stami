package utils

import (
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

func SingleRecipeDFS(recipe *scraper.Recipe, start string) scraper.TreeNode {
	root := scraper.TreeNode{Name: start}
	// implement here
	return root	
}

func MultipleRecipeDFS(recipe *scraper.Recipe, start string, numRecipe int) scraper.TreeNode {
	root := scraper.TreeNode{Name: start}
	// implement here
	return root	
}
