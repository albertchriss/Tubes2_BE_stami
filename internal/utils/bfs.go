package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/albertchriss/Tubes2_BE_stami/internal/app/socket"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

// type node struct {
// 	Id string `json:"id"`
// 	Name string `json:"name"`
// 	Children
// }

// type socketResponse struct {
// 	Type     string `json:"type"`
// 	RootNode string `json:"rootNode"`
// }

func SingleRecipeBFS(recipe *scraper.Recipe, start string, liveUpdate bool, wsManager *socket.ClientManager) scraper.TreeNode {
	root := scraper.TreeNode{Name: start}
	queue := []*scraper.TreeNode{&root} // tambahkan root ke queue
	visited := make(map[string]bool)
	if liveUpdate {
		wsManager.BroadcastNode(root)
	}
	for len(queue) > 0 {
		if liveUpdate {
			time.Sleep(1500 * time.Millisecond) // Tambahkan delay 100ms
		}

		currNode := queue[0]
		queue = queue[1:]

		if visited[currNode.Name] {
			continue
		}
		visited[currNode.Name] = true

		if scraper.IsBaseElement(currNode.Name) {
			continue
		}
		combinations, found := (*recipe)[currNode.Name]
		if !found || len(combinations) == 0 {
			fmt.Printf("Peringatan: Tidak ditemukan resep untuk elemen perantara '%s'.\n", currNode.Name)
			currNode.Children = nil // Pastikan tidak ada anak
			continue
		}

		var next *scraper.Combination = nil

		for _, combination := range combinations {
			if combination.First() != start && combination.Second() != start {
				next = &combination
				break
			}
		}

		if next == nil {
			fmt.Printf("Peringatan: Tidak ditemukan kombinasi yang valid untuk elemen '%s'.\n", currNode.Name)
			currNode.Children = nil // Pastikan tidak ada anak
			continue
		}

		first, second := next.First(), next.Second()
		node := &scraper.TreeNode{Name: "+"}
		node.Children = []scraper.TreeNode{
			{Name: first},
			{Name: second},
		}
		currNode.Children = append(currNode.Children, *node)
		if liveUpdate {
			wsManager.BroadcastNode(root)
		}
		queue = append(queue, &node.Children[0], &node.Children[1])
		// fmt.Print("Node: ", currNode.Name, " -> ", first, " + ", second, "\n")

	}

	return root
}

func MultipleRecipeBFS(recipe *scraper.Recipe, start string, numRecipe int, liveUpdate bool, wsManager *socket.ClientManager) scraper.TreeNode {

	// Buat node root untuk elemen target
	root := scraper.TreeNode{Name: start}
	if liveUpdate {
		wsManager.BroadcastNode(root)
	}

	queue := []*scraper.TreeNode{&root} // tambahkan root ke queue

	var mutex sync.Mutex
	var wg sync.WaitGroup

	visited := make(map[string]bool)
	currNum := 1

	for len(queue) > 0 {
		if liveUpdate {
			time.Sleep(1500 * time.Millisecond) // Tambahkan delay 100ms
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

				combinations, found := (*recipe)[currNode.Name]

				if !found || len(combinations) == 0 {
					fmt.Printf("Peringatan: Tidak ditemukan resep untuk elemen perantara '%s'.\n", currNode.Name)
					currNode.Children = nil // Pastikan tidak ada anak
					return
				}

				for i, combination := range combinations {
					if i > 0 {
						mutex.Lock()
						if currNum >= numRecipe {
							mutex.Unlock()
							break
						} else {
							currNum++
						}
						mutex.Unlock()
					}
					if combination.First() == start || combination.Second() == start {
						continue
					}
					first, second := combination.First(), combination.Second()
					node := &scraper.TreeNode{Name: "+"}
					node.Children = []scraper.TreeNode{
						{Name: first},
						{Name: second},
					}
					currNode.Children = append(currNode.Children, *node)

					if liveUpdate {
						mutex.Lock()
						wsManager.BroadcastNode(root)
						mutex.Unlock()
					}

					mutex.Lock()
					currentQueue = append(currentQueue, &node.Children[0], &node.Children[1])
					mutex.Unlock()
				}
			}(node)
		}
		wg.Wait()
		queue = currentQueue
	}

	// Kembalikan pohon resep yang sudah dibangun, dimulai dari root.
	return root
}
