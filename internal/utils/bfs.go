package utils

import (
	"sync"
	"time"

	"github.com/albertchriss/Tubes2_BE_stami/internal/app/socket"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

func SingleRecipeBFS(recipe *scraper.Recipe, tier *scraper.Tier, start string, liveUpdate bool, wsManager *socket.ClientManager) scraper.SearchResult {
	startTime := time.Now()
	nodeCount := 0
	id := 0
	root := scraper.TreeNode{Name: start, Id: id, ImageSrc: (*tier)[start].ImageSrc}
	nodeCount++
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
			{Name: first, Id: id, ImageSrc: (*tier)[first].ImageSrc},
			{Name: second, Id: id + 1, ImageSrc: (*tier)[second].ImageSrc},
		}
		id++
		currNode.Children = append(currNode.Children, *node)
		nodeCount += 2
		if liveUpdate {
			time.Sleep(500 * time.Millisecond)
			wsManager.BroadcastNode(root)
		}
		queue = append(queue, &node.Children[0], &node.Children[1])
	}

	duration := time.Since(startTime)
	return scraper.SearchResult{Tree: root, NodeCount: nodeCount, TimeTaken: duration.Nanoseconds()}
}

func MultipleRecipeBFS(recipe *scraper.Recipe, tier *scraper.Tier, start string, numRecipe int, liveUpdate bool, wsManager *socket.ClientManager) scraper.SearchResult {
	startTime := time.Now()
	nodeCount := 0
	id := 0
	// Buat node root untuk elemen target
	root := scraper.TreeNode{Name: start, Id: id, ImageSrc: (*tier)[start].ImageSrc}
	nodeCount++
	if liveUpdate {
		wsManager.BroadcastNode(root)
	}

	queue := []*scraper.TreeNode{&root} // tambahkan root ke queue

	var mutex sync.Mutex
	var wg sync.WaitGroup
	parent := map[int]*scraper.TreeNode{}
	numPath := map[int]int{}
	done := false

	for len(queue) > 0 {

		currentQueue := []*scraper.TreeNode{}

		for _, node := range queue {

			// jalankan go routine untuk setiap node pada queue saat ini
			wg.Add(1)
			go func(currNode *scraper.TreeNode) {
				defer wg.Done()
				if scraper.IsBaseElement(currNode.Name) {
					return
				}
				combinations := (*recipe)[currNode.Name]
				for i, combination := range combinations {
					mutex.Lock()
					currId := id
					id += 3
					mutex.Unlock()
					first, second := combination.First(), combination.Second()
					node := &scraper.TreeNode{Name: "+", Id: currId + 1}
					node.Children = []scraper.TreeNode{
						{Name: first, Id: currId + 2, ImageSrc: (*tier)[first].ImageSrc},
						{Name: second, Id: currId + 3, ImageSrc: (*tier)[second].ImageSrc},
					}

					mutex.Lock()
					node.InitParAndNum(currNode, &parent, &numPath)
					if done && i > 0 {
						mutex.Unlock()
						break
					}
					currNode.Children = append(currNode.Children, *node)
					currentQueue = append(currentQueue, &node.Children[0], &node.Children[1])
					if i > 0 {
						num := currNode.CountNumRecipe(&parent, &numPath)
						if num >= numRecipe {
							done = true
						}
					}
					mutex.Unlock()

					if liveUpdate {
						mutex.Lock()
						time.Sleep(500 * time.Millisecond) // Tambahkan delay 500ms
						wsManager.BroadcastNode(root)
						mutex.Unlock()
					}
				}
			}(node)
		}
		wg.Wait()

		mutex.Lock()
		nodeCount += len(currentQueue)
		mutex.Unlock()

		queue = currentQueue
	}

	duration := time.Since(startTime)
	return scraper.SearchResult{Tree: root, NodeCount: nodeCount, TimeTaken: duration.Nanoseconds()}
}
