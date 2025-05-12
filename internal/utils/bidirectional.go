package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
	"maps"
	"slices"
)

var globalNodeIdCounter int

func findNodeInTreeBFSInternal(root *scraper.TreeNode, name string) *scraper.TreeNode {
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
			childNode := &curr.Children[i]
			if childNode.Name == "+" {
				for j := range childNode.Children {
					q = append(q, &childNode.Children[j])
				}
			}
		}
	}
	return nil
}

func getBaseElementsInternal(tierMap *scraper.Tier) []string {
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

func reconstructSinglePathRecursive(
	currentNodeName string,
	recipeData *scraper.Recipe,
	stopCondition func(name string) bool,
	pathVisited map[string]bool,
	currentDepth int,
	maxRecursiveDepth int,
) scraper.TreeNode {

	globalNodeIdCounter++
	rootNode := scraper.TreeNode{Id: globalNodeIdCounter, Name: currentNodeName}

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

	combination := availableRecipes[0]
	ingredient1Name := combination.First()
	ingredient2Name := combination.Second()

	if ingredient1Name == "" && ingredient2Name == "" {
		return rootNode
	}

	var childrenForPlusNode []scraper.TreeNode
	childPathVisited1 := make(map[string]bool)
	for k, v := range pathVisited {
		childPathVisited1[k] = v
	}

	if ingredient1Name != "" {
		child1Tree := reconstructSinglePathRecursive(ingredient1Name, recipeData, stopCondition, childPathVisited1, currentDepth+1, maxRecursiveDepth)
		childrenForPlusNode = append(childrenForPlusNode, child1Tree)
	}

	childPathVisited2 := make(map[string]bool)
	maps.Copy(childPathVisited2, pathVisited)
	if ingredient2Name != "" {
		child2Tree := reconstructSinglePathRecursive(ingredient2Name, recipeData, stopCondition, childPathVisited2, currentDepth+1, maxRecursiveDepth)
		childrenForPlusNode = append(childrenForPlusNode, child2Tree)
	}

	if len(childrenForPlusNode) > 0 {
		globalNodeIdCounter++
		plusNode := scraper.TreeNode{Id: globalNodeIdCounter, Name: "+", Children: childrenForPlusNode}
		rootNode.Children = append(rootNode.Children, plusNode)
	}
	return rootNode
}

func solveSingleElementBidirectionally(
	targetElementName string,
	recipe *scraper.Recipe,
	tierMap *scraper.Tier,
) scraper.TreeNode {

	if recipe == nil {
		globalNodeIdCounter++
		return scraper.TreeNode{Id: globalNodeIdCounter, Name: "Error: Recipe data is nil."}
	}
	if scraper.IsBaseElement(targetElementName) {
		globalNodeIdCounter++
		return scraper.TreeNode{Id: globalNodeIdCounter, Name: targetElementName}
	}

	recipesForTarget, targetExists := (*recipe)[targetElementName]
	if !targetExists || len(recipesForTarget) == 0 {
		globalNodeIdCounter++
		return scraper.TreeNode{Id: globalNodeIdCounter, Name: fmt.Sprintf("Element '%s' no recipes.", targetElementName)}
	}

	visitedFwd := make(map[string]bool)
	baseElementsList := getBaseElementsInternal(tierMap)

	if len(baseElementsList) == 0 {
		globalNodeIdCounter++
		return scraper.TreeNode{Id: globalNodeIdCounter, Name: "Error: No base elements."}
	}

	currentFwdQueue := make([]string, 0, len(baseElementsList))
	for _, el := range baseElementsList {
		if !visitedFwd[el] {
			visitedFwd[el] = true
			currentFwdQueue = append(currentFwdQueue, el)
		}
	}

	visitedBwd := make(map[string]bool)
	currentBwdQueue := []string{targetElementName}
	visitedBwd[targetElementName] = true

	collectedMeetingNodes := make(map[string]bool)
	orderedMeetingNodes := []string{}

	const maxSearchDepth = 20
	const maxReconstructionDepth = 15

	for depth := range maxSearchDepth {
		if len(orderedMeetingNodes) > 0 && depth > 2 {
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
					canBeFormed := (ing1 == ing2 && visitedFwd[ing1]) || (visitedFwd[ing1] && visitedFwd[ing2])

					if canBeFormed {
						foundInNext := slices.Contains(nextFwdQueue, product)
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

			for _, ing := range []string{chosenCombBwd.First(), chosenCombBwd.Second()} {
				if ing != "" && !visitedBwd[ing] {
					visitedBwd[ing] = true
					foundInNextBwd := slices.Contains(nextBwdQueue, ing)
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
		stopConditionStandard := func(name string) bool { return scraper.IsBaseElement(name) }
		pathVisitedStandard := make(map[string]bool)
		return reconstructSinglePathRecursive(targetElementName, recipe, stopConditionStandard, pathVisitedStandard, 0, maxReconstructionDepth)
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
		return name == chosenMeetingNode || (scraper.IsBaseElement(name) && name != targetElementName)
	}
	pathVisitedSeg1 := make(map[string]bool)
	segmentStartToMeeting := reconstructSinglePathRecursive(targetElementName, recipe, stopConditionForSeg1, pathVisitedSeg1, 0, maxReconstructionDepth)

	var segmentMeetingToBase scraper.TreeNode
	if scraper.IsBaseElement(chosenMeetingNode) {
		globalNodeIdCounter++
		segmentMeetingToBase = scraper.TreeNode{}
	} else {
		stopConditionForSeg2 := func(name string) bool { return scraper.IsBaseElement(name) }
		pathVisitedSeg2 := make(map[string]bool)
		segmentMeetingToBase = reconstructSinglePathRecursive(chosenMeetingNode, recipe, stopConditionForSeg2, pathVisitedSeg2, 0, maxReconstructionDepth)
	}

	if segmentStartToMeeting.Name == chosenMeetingNode {
		if len(segmentMeetingToBase.Children) > 0 {
			segmentStartToMeeting.Children = segmentMeetingToBase.Children
		} 
		return segmentStartToMeeting
	}

	nodeToAttach := findNodeInTreeBFSInternal(&segmentStartToMeeting, chosenMeetingNode)

	if nodeToAttach != nil {
		if len(segmentMeetingToBase.Children) > 0 {
			nodeToAttach.Children = segmentMeetingToBase.Children
		} 
	}
	return segmentStartToMeeting
}

func BidirectionalSearch(
	recipe *scraper.Recipe,
	tierMap *scraper.Tier,
	startElementName string,
	numRecipesToFind int,
) scraper.TreeNode {
	globalNodeIdCounter = 0

	rootNode := scraper.TreeNode{Id: globalNodeIdCounter, Name: startElementName}

	if scraper.IsBaseElement(startElementName) {
		return rootNode
	}

	initialCombinations, exists := (*recipe)[startElementName]
	if !exists || len(initialCombinations) == 0 {
		rootNode.Name = fmt.Sprintf("No recipes for %s", startElementName)
		return rootNode
	}

	numActualRecipes := len(initialCombinations)
	if numRecipesToFind > 0 && numRecipesToFind < numActualRecipes {
		numActualRecipes = numRecipesToFind
	}
	if numRecipesToFind < 1 {
		numActualRecipes = 1
		if len(initialCombinations) == 0 {
			numActualRecipes = 0
		}
	}

	for i := 0; i < numActualRecipes; i++ {
		combination := initialCombinations[i]
		ing1 := combination.First()
		ing2 := combination.Second()

		var treeIngredient1, treeIngredient2 scraper.TreeNode

		if ing1 != "" {
			treeIngredient1 = solveSingleElementBidirectionally(ing1, recipe, tierMap)
		}
		if ing2 != "" {
			treeIngredient2 = solveSingleElementBidirectionally(ing2, recipe, tierMap)
		}

		globalNodeIdCounter++
		plusNode := scraper.TreeNode{
			Id:   globalNodeIdCounter,
			Name: "+",
		}

		isError1 := false
		if treeIngredient1.Name != "" {
			if len(treeIngredient1.Name) >= 5 && treeIngredient1.Name[:5] == "Error" {
				isError1 = true
			}
			if !isError1 && len(treeIngredient1.Name) >= 12 && strings.HasSuffix(treeIngredient1.Name, "no recipes.") {
				isError1 = true
			}
			if !isError1 && len(treeIngredient1.Name) >= 20 && strings.HasSuffix(treeIngredient1.Name, "no base elements.") {
				isError1 = true
			}
			 if !isError1 && len(treeIngredient1.Name) >= 25 && strings.HasSuffix(treeIngredient1.Name, "Recipe data is nil.") {
				isError1 = true
			}
		} else {
			isError1 = true
		}


		isError2 := false
		if treeIngredient2.Name != "" {
			if len(treeIngredient2.Name) >= 5 && treeIngredient2.Name[:5] == "Error" {
				isError2 = true
			}
			if !isError2 && len(treeIngredient2.Name) >= 12 && strings.HasSuffix(treeIngredient2.Name, "no recipes.") {
				isError2 = true
			}
			if !isError2 && len(treeIngredient2.Name) >= 20 && strings.HasSuffix(treeIngredient2.Name, "no base elements.") {
				isError2 = true
			}
			if !isError2 && len(treeIngredient2.Name) >= 25 && strings.HasSuffix(treeIngredient2.Name, "Recipe data is nil.") {
				isError2 = true
			}
		} else {
			isError2 = true
		}


		if !isError1 {
			plusNode.Children = append(plusNode.Children, treeIngredient1)
		}
		if !isError2 {
			plusNode.Children = append(plusNode.Children, treeIngredient2)
		}

		if len(plusNode.Children) > 0 {
			rootNode.Children = append(rootNode.Children, plusNode)
		}
	}
	return rootNode
}