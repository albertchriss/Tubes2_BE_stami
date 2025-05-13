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
	Parent     *map[int]*scraper.TreeNode
	NumPath    *map[int]int
	Mutex      *sync.Mutex
	Wg         *sync.WaitGroup
	Done       *bool
	LiveUpdate bool
	WsManager  *socket.ClientManager
}

type SingleHelperParams struct {
	Root       *scraper.TreeNode
	Id         *int
	LiveUpdate bool
	WsManager  *socket.ClientManager
}

func SingleRecipeDFS(recipe *scraper.Recipe, tier *scraper.Tier, start string, liveUpdate bool, wsManager *socket.ClientManager) scraper.TreeNode {
	id := 0
	root := scraper.TreeNode{Name: start, Id: id, ImageSrc: (*tier)[start].ImageSrc}
	if liveUpdate {
		wsManager.BroadcastNode(root)
	}
	params := &SingleHelperParams{
		Root:       &root,
		Id:         &id,
		LiveUpdate: liveUpdate,
		WsManager:  wsManager,
	}
	SingleDFSHelper(recipe, tier, start, params, &root)
	return root
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

	if liveUpdate {
		time.Sleep(500 * time.Millisecond)
		wsManager.BroadcastNode(*root)
	}

	SingleDFSHelper(recipe, tier, first, params, &node.Children[0])
	SingleDFSHelper(recipe, tier, second, params, &node.Children[1])
}

func MultipleRecipeDFS(recipe *scraper.Recipe, tier *scraper.Tier, start string, numRecipe int, liveUpdate bool, wsManager *socket.ClientManager) scraper.TreeNode {
	id := 0
	done := false
	var mutex sync.Mutex
	var wg sync.WaitGroup

	root := scraper.TreeNode{Name: start, Id: id, ImageSrc: (*tier)[start].ImageSrc}

	if liveUpdate {
		wsManager.BroadcastNode(root)
	}

	wg.Add(1)
	params := &MultHelperParams{
		Root:       &root,
		Id:         &id,
		Parent:     &map[int]*scraper.TreeNode{},
		NumPath:    &map[int]int{},
		Done:       &done,
		Mutex:      &mutex,
		Wg:         &wg,
		LiveUpdate: liveUpdate,
		WsManager:  wsManager,
	}
	go MultipleRecipeHelper(recipe, tier, start, numRecipe, params, &root)
	wg.Wait()
	return root
}

func MultipleRecipeHelper(recipe *scraper.Recipe, tier *scraper.Tier, name string, numRecipe int, params *MultHelperParams, currNode *scraper.TreeNode) {
	wg := params.Wg
	mutex := params.Mutex

	mutex.Lock()
	parent := params.Parent
	numPath := params.NumPath
	done := params.Done
	liveUpdate := params.LiveUpdate
	wsManager := params.WsManager
	root := params.Root
	id := params.Id
	mutex.Unlock()

	defer wg.Done()

	if scraper.IsBaseElement(name) {
		return
	}
	combinations := (*recipe)[name]

	for i, combination := range combinations {

		mutex.Lock()
		currId := (*id)
		(*id) += 3
		mutex.Unlock()

		first, second := combination.First(), combination.Second()
		node := &scraper.TreeNode{Name: "+", Id: currId + 1}
		node.Children = []scraper.TreeNode{
			{Name: first, Id: currId + 2, ImageSrc: (*tier)[first].ImageSrc},
			{Name: second, Id: currId + 3, ImageSrc: (*tier)[second].ImageSrc},
		}

		mutex.Lock()
		node.InitParAndNum(currNode, parent, numPath)
		if (*done) && i > 0 {
			mutex.Unlock()
			break
		}
		currNode.Children = append(currNode.Children, *node)
		if i > 0 {
			num := currNode.CountNumRecipe(parent, numPath)
			if num >= numRecipe {
				(*done) = true
			}
		}
		mutex.Unlock()
		wg.Add(1)
		go MultipleRecipeHelper(recipe, tier, first, numRecipe, params, &node.Children[0])
		wg.Add(1)
		go MultipleRecipeHelper(recipe, tier, second, numRecipe, params, &node.Children[1])

		if liveUpdate {
			mutex.Lock()
			time.Sleep(300 * time.Millisecond)
			wsManager.BroadcastNode(*root)
			mutex.Unlock()
		}

	}
	// wg.Wait()
}
