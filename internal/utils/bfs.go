package utils

import (
	"sync"
	"time"

	"github.com/albertchriss/Tubes2_BE_stami/internal/app/socket"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

func SingleRecipeBFS(recipe *scraper.Recipe, start string, liveUpdate bool, wsManager *socket.ClientManager) scraper.TreeNode {
	id := 0
	root := scraper.TreeNode{Name: start, Id: id}
	if liveUpdate {
		wsManager.BroadcastNode(root)
	}
	queue := []*scraper.TreeNode{&root}
	for len(queue) > 0 {

		currNode := queue[0]
		queue = queue[1:]

		if scraper.IsBaseElement(currNode.Name) {
			continue
		}
		combinations := (*recipe)[currNode.Name]
		next := combinations[0]
		first, second := next.First(), next.Second()
		id++
		node := &scraper.TreeNode{Name: "+", Id: id}
		id++
		node.Children = []scraper.TreeNode{
			{Name: first, Id: id},
			{Name: second, Id: id + 1},
		}
		id++
		currNode.Children = append(currNode.Children, *node)
		if liveUpdate {
			time.Sleep(500 * time.Millisecond)
			wsManager.BroadcastNode(root)
		}
		queue = append(queue, &node.Children[0], &node.Children[1])
	}

	return root
}

func MultipleRecipeBFS(recipe *scraper.Recipe, start string, numRecipe int, liveUpdate bool, wsManager *socket.ClientManager) scraper.TreeNode {
	id := 0
	// Buat node root untuk elemen target
	root := scraper.TreeNode{Name: start, Id: id}
	if liveUpdate {
		wsManager.BroadcastNode(root)
	}

	queue := []*scraper.TreeNode{&root} // tambahkan root ke queue

	var mutex sync.Mutex
	var wg sync.WaitGroup

	currNum := 1

	for len(queue) > 0 {

		currentQueue := []*scraper.TreeNode{}

		for _, node := range queue {
			wg.Add(1)
			go func(currNode *scraper.TreeNode) {
				defer wg.Done()
				if scraper.IsBaseElement(currNode.Name) {
					return
				}
				combinations := (*recipe)[currNode.Name]
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
					first, second := combination.First(), combination.Second()
					mutex.Lock()
					id++
					node := &scraper.TreeNode{Name: "+", Id: id}
					id++
					node.Children = []scraper.TreeNode{
						{Name: first, Id: id},
						{Name: second, Id: id + 1},
					}
					id++
					mutex.Unlock()
					currNode.Children = append(currNode.Children, *node)

					if liveUpdate {
						time.Sleep(500 * time.Millisecond) // Tambahkan delay 100ms
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

	return root
}
