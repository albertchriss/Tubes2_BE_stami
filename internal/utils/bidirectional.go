package utils

import (
	"fmt"
	"sort"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

func findNodeInTreeBFS(root *scraper.TreeNode, name string) *scraper.TreeNode {
	if root == nil {
		return nil
	}
	q := []*scraper.TreeNode{root}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr.Name == name {
			return curr
		}

		for i := range curr.Children {
			plusNode := &curr.Children[i]
			if plusNode.Name == "+" {
				for j := range plusNode.Children {
					q = append(q, &plusNode.Children[j])
				}
			}
		}
	}
	return nil
}

func getBaseElements(tierMap *scraper.Tier) []string {
	var baseElements []string
	knownBaseElements := []string{"Air", "Earth", "Fire", "Water", "Time"}

	if tierMap != nil && len(*tierMap) > 0 {
		tempBaseElementsMap := make(map[string]bool)
		for el, tierVal := range *tierMap {
			if tierVal == 0 {
				tempBaseElementsMap[el] = true
			}
		}

		if len(tempBaseElementsMap) > 0 {
			for el := range tempBaseElementsMap {
				baseElements = append(baseElements, el)
			}
			sort.Strings(baseElements)
			return baseElements
		}
	}
	sort.Strings(knownBaseElements)
	return knownBaseElements
}

func reconstructPathMultipleRecipesInternal(
	currentNodeName string,
	recipeData *scraper.Recipe,
	stopCondition func(name string) bool,
	maxRecipesPerNodeToDisplay int,
	pathVisited map[string]bool,
	currentDepth int,
	maxRecursiveDepth int,
) scraper.TreeNode {

	rootNode := scraper.TreeNode{Name: currentNodeName}

	if currentDepth > maxRecursiveDepth {
		return rootNode
	}

	if stopCondition(currentNodeName) {
		return rootNode
	}

	if pathVisited[currentNodeName] {
		return rootNode
	}
	pathVisited[currentNodeName] = true
	defer delete(pathVisited, currentNodeName)

	availableRecipes, recipesExist := (*recipeData)[currentNodeName]
	if !recipesExist || len(availableRecipes) == 0 {
		return rootNode
	}

	recipesToProcess := len(availableRecipes)
	if recipesToProcess > maxRecipesPerNodeToDisplay {
		recipesToProcess = maxRecipesPerNodeToDisplay
	}

	for i := 0; i < recipesToProcess; i++ {
		combination := availableRecipes[i]
		ingredient1Name := combination.First()
		ingredient2Name := combination.Second()

		if ingredient1Name == "" && ingredient2Name == "" {
			continue
		}

		var childrenForPlusNode []scraper.TreeNode
		if ingredient1Name != "" {
			child1Tree := reconstructPathMultipleRecipesInternal(ingredient1Name, recipeData, stopCondition, maxRecipesPerNodeToDisplay, pathVisited, currentDepth+1, maxRecursiveDepth)
			childrenForPlusNode = append(childrenForPlusNode, child1Tree)
		}
		if ingredient2Name != "" {
			child2Tree := reconstructPathMultipleRecipesInternal(ingredient2Name, recipeData, stopCondition, maxRecipesPerNodeToDisplay, pathVisited, currentDepth+1, maxRecursiveDepth)
			childrenForPlusNode = append(childrenForPlusNode, child2Tree)
		}

		if len(childrenForPlusNode) > 0 {
			plusNode := scraper.TreeNode{Name: "+", Children: childrenForPlusNode}
			rootNode.Children = append(rootNode.Children, plusNode)
		}
	}
	return rootNode
}

func BidirectionalSearch(
	recipe *scraper.Recipe,
	tierMap *scraper.Tier,
	startElementName string,
	maxRecipesPerNodeInTree int,
) scraper.TreeNode {

	if recipe == nil {
		return scraper.TreeNode{Name: "Error: Recipe data is nil."}
	}
	if scraper.IsBaseElement(startElementName) {
		return scraper.TreeNode{Name: startElementName}
	}
	if _, exists := (*recipe)[startElementName]; !exists {
		recipesForStart, hasRecipes := (*recipe)[startElementName]
		if !hasRecipes || len(recipesForStart) == 0 {
			return scraper.TreeNode{Name: fmt.Sprintf("Target element '%s' not found in recipes or has no defined combinations.", startElementName)}
		}
	}

	if maxRecipesPerNodeInTree < 1 {
		maxRecipesPerNodeInTree = 1
	}

	visitedFwd := make(map[string]bool)
	parentsFwd := make(map[string]scraper.Combination) 
	baseElementsList := getBaseElements(tierMap)

	if len(baseElementsList) == 0 {
		return scraper.TreeNode{Name: "Error: No base elements found."}
	}

	currentFwdQueue := make([]string, 0, len(baseElementsList))
	for _, el := range baseElementsList {
		if !visitedFwd[el] {
			visitedFwd[el] = true
			currentFwdQueue = append(currentFwdQueue, el)
		}
	}

	visitedBwd := make(map[string]bool)
	parentsBwd := make(map[string]scraper.Combination) 
	currentBwdQueue := []string{startElementName}
	visitedBwd[startElementName] = true

	collectedMeetingNodes := make(map[string]bool)
	orderedMeetingNodes := []string{}

	const maxSearchDepth = 20
	const maxReconstructionDepth = 15

	for depth := 0; depth < maxSearchDepth; depth++ {
		if len(orderedMeetingNodes) > 0 && depth > 1 {
			if len(orderedMeetingNodes) > 5 {
				break
			}
		}

		nextFwdQueue := []string{}
		for _, nodeFromBase := range currentFwdQueue {
			if visitedBwd[nodeFromBase] && !collectedMeetingNodes[nodeFromBase] {
				collectedMeetingNodes[nodeFromBase] = true
				orderedMeetingNodes = append(orderedMeetingNodes, nodeFromBase)
			}

			for product, combinations := range *recipe {
				if visitedFwd[product] {
					continue
				}
				for _, comb := range combinations {
					ing1, ing2 := comb.First(), comb.Second()
					canBeFormed := false
					if ing1 == ing2 {
						if visitedFwd[ing1] {
							canBeFormed = true
						}
					} else {
						if visitedFwd[ing1] && visitedFwd[ing2] {
							canBeFormed = true
						}
					}

					if canBeFormed {
						if _, ok := parentsFwd[product]; !ok {
							parentsFwd[product] = comb
						}
						foundInNext := false
						for _, item := range nextFwdQueue {
							if item == product {
								foundInNext = true
								break
							}
						}
						if !foundInNext {
							nextFwdQueue = append(nextFwdQueue, product)
						}
						break
					}
				}
			}
		}
		for _, n := range nextFwdQueue {
			visitedFwd[n] = true
		}
		currentFwdQueue = nextFwdQueue

		nextBwdQueue := []string{}
		for _, nodeToMake := range currentBwdQueue {
			if visitedFwd[nodeToMake] && !collectedMeetingNodes[nodeToMake] {
				collectedMeetingNodes[nodeToMake] = true
				orderedMeetingNodes = append(orderedMeetingNodes, nodeToMake)
			}

			if scraper.IsBaseElement(nodeToMake) {
				continue
			}

			recipesForNode, exists := (*recipe)[nodeToMake]
			if !exists || len(recipesForNode) == 0 {
				continue
			}

			chosenCombBwd := recipesForNode[0]
			if _, ok := parentsBwd[nodeToMake]; !ok {
				parentsBwd[nodeToMake] = chosenCombBwd
			}

			for _, ing := range []string{chosenCombBwd.First(), chosenCombBwd.Second()} {
				if ing != "" && !visitedBwd[ing] {
					visitedBwd[ing] = true
					foundInNextBwd := false
					for _, item := range nextBwdQueue {
						if item == ing {
							foundInNextBwd = true
							break
						}
					}
					if !foundInNextBwd {
						nextBwdQueue = append(nextBwdQueue, ing)
					}
				}
			}
		}
		currentBwdQueue = nextBwdQueue

		if len(currentFwdQueue) == 0 && len(currentBwdQueue) == 0 && depth > 1 {
			break
		}
	}

	if len(orderedMeetingNodes) == 0 {
		stopConditionStandard := func(name string) bool {
			return scraper.IsBaseElement(name)
		}
		pathVisitedStandard := make(map[string]bool)
		return reconstructPathMultipleRecipesInternal(startElementName, recipe, stopConditionStandard, maxRecipesPerNodeInTree, pathVisitedStandard, 0, maxReconstructionDepth)
	}

	sort.Slice(orderedMeetingNodes, func(i, j int) bool {
		tierI, okI := (*tierMap)[orderedMeetingNodes[i]]
		tierJ, okJ := (*tierMap)[orderedMeetingNodes[j]]
		if !okI {
			tierI = maxSearchDepth + 1
		}
		if !okJ {
			tierJ = maxSearchDepth + 1
		}
		if tierI != tierJ {
			return tierI < tierJ
		}
		return orderedMeetingNodes[i] < orderedMeetingNodes[j]
	})

	chosenMeetingNode := orderedMeetingNodes[0]

	stopConditionForSeg1 := func(name string) bool {
		if name == chosenMeetingNode && name != startElementName {
			return true
		}
		return scraper.IsBaseElement(name) && name != startElementName
	}
	pathVisitedSeg1 := make(map[string]bool)
	segmentStartToMeeting := reconstructPathMultipleRecipesInternal(startElementName, recipe, stopConditionForSeg1, maxRecipesPerNodeInTree, pathVisitedSeg1, 0, maxReconstructionDepth)

	var segmentMeetingToBase scraper.TreeNode
	if scraper.IsBaseElement(chosenMeetingNode) {
		segmentMeetingToBase = scraper.TreeNode{}
	} else {
		stopConditionForSeg2 := func(name string) bool {
			return scraper.IsBaseElement(name)
		}
		pathVisitedSeg2 := make(map[string]bool)
		segmentMeetingToBase = reconstructPathMultipleRecipesInternal(chosenMeetingNode, recipe, stopConditionForSeg2, maxRecipesPerNodeInTree, pathVisitedSeg2, 0, maxReconstructionDepth)
	}

	if startElementName == chosenMeetingNode {
		segmentStartToMeeting.Children = segmentMeetingToBase.Children
		return segmentStartToMeeting
	}

	nodeToAttach := findNodeInTreeBFS(&segmentStartToMeeting, chosenMeetingNode)
	if nodeToAttach != nil {
		nodeToAttach.Children = segmentMeetingToBase.Children
	}

	return segmentStartToMeeting
}