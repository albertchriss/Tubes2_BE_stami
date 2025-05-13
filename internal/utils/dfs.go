package utils

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/albertchriss/Tubes2_BE_stami/internal/app/socket"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

type MultHelperParams struct {
	Root       *scraper.TreeNode
	Id         *int
	Count      *int
	NodeCount  *int32
	Mutex      *sync.Mutex
	Wg         *sync.WaitGroup
	LiveUpdate bool
	WsManager  *socket.ClientManager
}

type SingleHelperParams struct {
	Root       *scraper.TreeNode
	Id         *int
	NodeCount  *int
	LiveUpdate bool
	WsManager  *socket.ClientManager
}

func SingleRecipeDFS(recipe *scraper.Recipe, tier *scraper.Tier, start string, liveUpdate bool, wsManager *socket.ClientManager) scraper.SearchResult {
	startTime := time.Now()
	nodeCount := 0
	id := 0
	root := scraper.TreeNode{Name: start, Id: id, ImageSrc: (*tier)[start].ImageSrc}
	nodeCount++
	if liveUpdate {
		wsManager.BroadcastNode(root)
	}
	params := &SingleHelperParams{
		Root:       &root,
		Id:         &id,
		NodeCount:  &nodeCount,
		LiveUpdate: liveUpdate,
		WsManager:  wsManager,
	}
	SingleDFSHelper(recipe, tier, start, params, &root)
	duration := time.Since(startTime)
	return scraper.SearchResult{Tree: root, NodeCount: nodeCount, TimeTaken: duration.Nanoseconds()}
}

func SingleDFSHelper(recipe *scraper.Recipe, tier *scraper.Tier, start string, params *SingleHelperParams, currNode *scraper.TreeNode) {

	root := params.Root
	liveUpdate := params.LiveUpdate
	wsManager := params.WsManager
	id := params.Id

	if scraper.IsBaseElement(start) {
		return
	}

	combinations := (*recipe)[start]

	next := combinations[0]
	first, second := next.First(), next.Second()

	(*id)++
	node := &scraper.TreeNode{Name: "+", Id: (*id)}
	(*id)++
	node.Children = []scraper.TreeNode{
		{Name: first, Id: (*id), ImageSrc: (*tier)[first].ImageSrc},
		{Name: second, Id: (*id) + 1, ImageSrc: (*tier)[second].ImageSrc},
	}
	(*id)++
	currNode.Children = append(currNode.Children, *node)

	(*params.NodeCount) += 2
	
	if liveUpdate {
		time.Sleep(500 * time.Millisecond)
		wsManager.BroadcastNode(*root)
	}

	SingleDFSHelper(recipe, tier, first, params, &node.Children[0])
	SingleDFSHelper(recipe, tier, second, params, &node.Children[1])
}

func MultipleRecipeDFS(recipe *scraper.Recipe, tier *scraper.Tier, start string, numRecipe int, liveUpdate bool, wsManager *socket.ClientManager) scraper.SearchResult {
	startTime := time.Now()
	count := 1
	id := 0
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var nodeCounter int32 = 0

	root := scraper.TreeNode{Name: start, Id: id, ImageSrc: (*tier)[start].ImageSrc}

	if liveUpdate {
		wsManager.BroadcastNode(root)
	}

	wg.Add(1)
	params := &MultHelperParams{
		Root:       &root,
		Id:         &id,
		Count:      &count,
		NodeCount:  &nodeCounter,
		Mutex:      &mutex,
		Wg:         &wg,
		LiveUpdate: liveUpdate,
		WsManager:  wsManager,
	}
	go MultipleRecipeHelper(recipe, tier, start, numRecipe, params, &root)
	wg.Wait()
	duration := time.Since(startTime)
	finalNodeCount := int(atomic.LoadInt32(&nodeCounter))
	return scraper.SearchResult{Tree: root, NodeCount: finalNodeCount, TimeTaken: duration.Nanoseconds()}
}

func MultipleRecipeHelper(recipe *scraper.Recipe, tier *scraper.Tier, name string, numRecipe int, params *MultHelperParams, currNode *scraper.TreeNode) {
	count := params.Count
	wg := params.Wg
	mutex := params.Mutex
	liveUpdate := params.LiveUpdate
	wsManager := params.WsManager
	root := params.Root
	id := params.Id

	defer wg.Done()

	if scraper.IsBaseElement(name) {
		return
	}
	combinations := (*recipe)[name]

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
		mutex.Lock()
		(*id)++
		node := &scraper.TreeNode{Name: "+", Id: *id}
		(*id)++
		node.Children = []scraper.TreeNode{
			{Name: first, Id: *id, ImageSrc: (*tier)[first].ImageSrc},
			{Name: second, Id: (*id) + 1, ImageSrc: (*tier)[second].ImageSrc},
		}
		(*id)++
		atomic.AddInt32(params.NodeCount, 2)
		mutex.Unlock()
		currNode.Children = append(currNode.Children, *node)
		if liveUpdate {
			time.Sleep(500 * time.Millisecond)
			wsManager.BroadcastNode(*root)
		}

		wg.Add(1)
		go MultipleRecipeHelper(recipe, tier, first, numRecipe, params, &node.Children[0])
		wg.Add(1)
		go MultipleRecipeHelper(recipe, tier, second, numRecipe, params, &node.Children[1])
	}
}
