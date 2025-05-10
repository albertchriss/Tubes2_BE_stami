package utils

import (
	"fmt"
	"sync"

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
	count := 1
	visited := make(map[string]bool)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	root := scraper.TreeNode{Name: start}
	wg.Add(1)
	go MultipleRecipeHelper(recipe, &root, start, &count, numRecipe, visited, &mutex, &wg)
	wg.Wait()
	return root
}

func MultipleRecipeHelper(recipe *scraper.Recipe, root *scraper.TreeNode, name string, count *int, numRecipe int, visited map[string]bool, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	mutex.Lock()
	if visited[name] {
		// mutex.Unlock()
		return
	}
	visited[name] = true
	mutex.Unlock()

	if scraper.IsBaseElement(name) {
		return
	}

	combinations := (*recipe)[name]
	if len(combinations) == 0 {
		return
	}

	for i, combination := range combinations {
		if i > 0 {
			mutex.Lock()
			if *count >= numRecipe {
				mutex.Unlock()
				break
			}
			(*count)++
			mutex.Unlock()
		}

		first, second := combination.First(), combination.Second()
		node := scraper.TreeNode{Name: "+"}
		child1 := scraper.TreeNode{Name: first}
		child2 := scraper.TreeNode{Name: second}

		node.Children = []scraper.TreeNode{child1, child2}
		mutex.Lock()
		root.Children = append(root.Children, node)
		mutex.Unlock()

		if !scraper.IsBaseElement(first) {
			wg.Add(1)
			go MultipleRecipeHelper(recipe, &node.Children[0], first, count, numRecipe, copyVisited(visited, name), mutex, wg)
		}
		if !scraper.IsBaseElement(second) {
			wg.Add(1)
			go MultipleRecipeHelper(recipe, &node.Children[1], second, count, numRecipe, copyVisited(visited, name), mutex, wg)
		}
	}
}

func copyVisited(current map[string]bool, name string) map[string]bool {
	newVisited := make(map[string]bool)
	for k, v := range current {
		newVisited[k] = v
	}
	newVisited[name] = true
	return newVisited
}

// func MultipleRecipeDFS(recipe *scraper.Recipe, start string, numRecipe int) scraper.TreeNode {
// 	count := 1
// 	visited := make(map[string]bool)
// 	return MultipleRecipeHelper(recipe, start, &count, numRecipe, visited)
// }

// func MultipleRecipeHelper(recipe *scraper.Recipe, name string, count *int, numRecipe int, visited map[string]bool) scraper.TreeNode {
// 	root := scraper.TreeNode{Name: name}

// 	if visited[name] {
// 		return root
// 	}
// 	visited[name] = true

// 	if scraper.IsBaseElement(name) {
// 		return root
// 	}

// 	combinations := (*recipe)[name]
// 	if len(combinations) == 0 {
// 		return root
// 	}

// 	for i, combination := range combinations {
// 		if i > 0 {
// 			if *count >= numRecipe {
// 				break
// 			}
// 			(*count)++
// 		}

// 		first := combination.First()
// 		second := combination.Second()
// 		node := scraper.TreeNode{Name: "+"}

// 		var left, right scraper.TreeNode
// 		if scraper.IsBaseElement(first) {
// 			left = scraper.TreeNode{Name: first}
// 		} else {
// 			left = MultipleRecipeHelper(recipe, first, count, numRecipe, visited)
// 		}

// 		if scraper.IsBaseElement(second) {
// 			right = scraper.TreeNode{Name: second}
// 		} else {
// 			right = MultipleRecipeHelper(recipe, second, count, numRecipe, visited)
// 		}

// 		node.Children = []scraper.TreeNode{left, right}
// 		root.Children = append(root.Children, node)
// 	}

// 	return root
// }
