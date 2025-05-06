package utils

import (
	"fmt"
	"sync"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

// func SingleRecipeBFS(recipe *scraper.Recipe, start string) {
// 	var queue []string
// 	status := make(map[string]int)
// 	parent := make(map[string]int)
// 	nodeId := make(map[int]scraper.Combination)
// 	queue = append(queue, start)
// 	visited[start] = true
// 	for len(queue) > 0 {
// 		current := queue[0]
// 		queue = queue[1:]

// 		if scraper.IsBaseElement(current) {
// 			continue
// 		}

// 		for _, combination := range (*recipe)[current] {
// 			first, second := combination.First(), combination.Second()
// 			if !visited[first] {
// 				queue = append(queue, first)
// 				visited[first] = true
// 			}
// 			if !visited[second] {
// 				queue = append(queue, second)
// 				visited[second] = true
// 			}
// 		}
// 	}
// }

func MultipleRecipeBFS(recipe scraper.Recipe, start string) scraper.TreeNode {

	// Buat node root untuk elemen target
	root := scraper.TreeNode{Name: start}

	queue := []*scraper.TreeNode{&root} // tambahkan root ke queue

	var mutex sync.Mutex
	var wg sync.WaitGroup

	visited := make(map[string]bool)
	depth := 0

	for len(queue) > 0 {
		if depth > 3 {
			fmt.Println("Peringatan: Terlalu dalam, menghentikan pencarian.")
			break
		}

		currentQueue := []*scraper.TreeNode{}

		for _, node := range queue {
			wg.Add(1)
			go func(currNode *scraper.TreeNode) {
				defer wg.Done()

				// Cek apakah node sudah dikunjungi
				mutex.Lock()
				if visited[currNode.Name] {
					mutex.Unlock()
					return
				}
				visited[currNode.Name] = true
				mutex.Unlock()

				if scraper.IsBaseElement(currNode.Name) {
					return
				}

				combinations, found := recipe[currNode.Name]

				if !found || len(combinations) == 0 {
					fmt.Printf("Peringatan: Tidak ditemukan resep untuk elemen perantara '%s'.\n", currNode.Name)
					currNode.Children = nil // Pastikan tidak ada anak
					return
				}

				for i, combination := range combinations {
					if i > 1 {
						break
					}
					first, second := combination.First(), combination.Second()
					node := &scraper.TreeNode{Name: "+"}
					node.Children = []scraper.TreeNode{
						{Name: first},
						{Name: second},
					}
					currNode.Children = append(currNode.Children, *node)
					mutex.Lock()
					currentQueue = append(currentQueue, &node.Children[0], &node.Children[1])
					mutex.Unlock()
				}
			}(node)
		}
		wg.Wait()
		queue = currentQueue
		depth++
	}

	// Kembalikan pohon resep yang sudah dibangun, dimulai dari root.
	return root
}
