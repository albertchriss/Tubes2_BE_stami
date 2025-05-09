package utils

import (
	"fmt"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

func SingleRecipeDFS(recipe *scraper.Recipe, start string) scraper.TreeNode {
	visited := make(map[string]bool)
	return SingleDFSHelper(recipe, start, visited)
}

func SingleDFSHelper(recipe *scraper.Recipe, start string, visited map[string]bool) scraper.TreeNode {
	root := scraper.TreeNode{Name: start}
	if visited[start] {
		return root
	}
	visited[start] = true

	if scraper.IsBaseElement(start) {
		return root
	}

	combinations, found := (*recipe)[start]
	if !found || len(combinations) == 0 {
		fmt.Printf("Peringatan: Tidak ditemukan resep untuk elemen perantara '%s'.\n", start)
		return root
	}

	next := combinations[0]
	first, second := next.First(), next.Second()

	node := scraper.TreeNode{Name: "+"}
	node.Children = []scraper.TreeNode{
		SingleDFSHelper(recipe, first, visited),
		SingleDFSHelper(recipe, second, visited),
	}

	root.Children = append(root.Children, node)
	return root
}

func MultipleRecipeDFS(recipe *scraper.Recipe, start string, numRecipe int) scraper.TreeNode {
	root := scraper.TreeNode{Name: start}
	// implement here
	return root
}
