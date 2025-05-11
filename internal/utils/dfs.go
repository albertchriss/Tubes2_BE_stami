package utils

import (
	"sync"
	"time"

	"github.com/albertchriss/Tubes2_BE_stami/internal/app/socket"
	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

type MultHelperParams struct {
	Root       *scraper.TreeNode
	Id         *int
	Count      *int
	Mutex      *sync.Mutex
	Wg         *sync.WaitGroup
	LiveUpdate bool
	WsManager  *socket.ClientManager
}

type SingleHelperParams struct {
	Root       *scraper.TreeNode
	Id         *int
	LiveUpdate bool
	WsManager  *socket.ClientManager
}

func SingleRecipeDFS(recipe *scraper.Recipe, start string, liveUpdate bool, wsManager *socket.ClientManager) scraper.TreeNode {
	id := 0
	root := scraper.TreeNode{Name: start}
	if liveUpdate {
		wsManager.BroadcastNode(root)
	}
	params := &SingleHelperParams{
		Root:       &root,
		Id:         &id,
		LiveUpdate: liveUpdate,
		WsManager:  wsManager,
	}
	SingleDFSHelper(recipe, start, params, &root)
	return root
}

func SingleDFSHelper(recipe *scraper.Recipe, start string, params *SingleHelperParams, currNode *scraper.TreeNode) {

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
		{Name: first, Id: (*id)},
		{Name: second, Id: (*id) + 1},
	}
	(*id)++
	currNode.Children = append(currNode.Children, *node)

	if liveUpdate {
		time.Sleep(500 * time.Millisecond)
		wsManager.BroadcastNode(*root)
	}

	SingleDFSHelper(recipe, first, params, &node.Children[0])
	SingleDFSHelper(recipe, second, params, &node.Children[1])
}

func MultipleRecipeDFS(recipe *scraper.Recipe, start string, numRecipe int, liveUpdate bool, wsManager *socket.ClientManager) scraper.TreeNode {
	count := 1
	id := 0
	var mutex sync.Mutex
	var wg sync.WaitGroup

	root := scraper.TreeNode{Name: start, Id: id}

	if liveUpdate {
		wsManager.BroadcastNode(root)
	}

	wg.Add(1)
	params := &MultHelperParams{
		Root:       &root,
		Id:         &id,
		Count:      &count,
		Mutex:      &mutex,
		Wg:         &wg,
		LiveUpdate: liveUpdate,
		WsManager:  wsManager,
	}
	go MultipleRecipeHelper(recipe, start, numRecipe, params, &root)
	wg.Wait()
	return root
}

func MultipleRecipeHelper(recipe *scraper.Recipe, name string, numRecipe int, params *MultHelperParams, currNode *scraper.TreeNode) {
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
			{Name: first, Id: *id},
			{Name: second, Id: (*id)+1},
		}
		(*id)++
		mutex.Unlock()
		currNode.Children = append(currNode.Children, *node)
		if liveUpdate {
			time.Sleep(500 * time.Millisecond)
			wsManager.BroadcastNode(*root)
		}

		wg.Add(1)
		go MultipleRecipeHelper(recipe, first, numRecipe, params, &node.Children[0])
		wg.Add(1)
		go MultipleRecipeHelper(recipe, second, numRecipe, params, &node.Children[1])
	}
}
