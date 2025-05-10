package utils

import (
	"fmt"
	"sort"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

func findNodeInTreeBFS(root *scraper.TreeNode, name string) *scraper.TreeNode {
	if root == nil { return nil}
	q := []*scraper.TreeNode{root}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr.Name == name {
			return curr
		}

		for i := range curr.Children {
			childRef := &curr.Children[i]
			if childRef.Name == "+" {
				for j := range childRef.Children {
					q = append(q, &childRef.Children[j])
				}
			} else {
				q = append(q, childRef)
			}
		}
	}
	return nil
}


func getBaseElements(tierMap *scraper.Tier) []string {
	var baseElements []string
	knownBaseElements := []string{"Air", "Earth", "Fire", "Water", "Time"}
	if tierMap != nil && len(*tierMap) > 0 {
		tempBaseElements := make(map[string]bool)
		for el, tier := range *tierMap {
			if tier == 0 {
				tempBaseElements[el] = true
			}
		}
		if len(tempBaseElements) > 0 {
			for el := range tempBaseElements {
				baseElements = append(baseElements, el)
			}
			return baseElements
		}
	}
	return knownBaseElements
}

func reconstructPathDeterministic(
	startNode string,
	parents map[string]scraper.Combination,
	recipeData *scraper.Recipe,
	stopAtNode string,
	isSegmentToMeetingNode bool) scraper.TreeNode {

	root := scraper.TreeNode{Name: startNode}
	queue := []*scraper.TreeNode{&root}
	expandedInThisReconstruction := make(map[string]bool)

	head := 0
	for head < len(queue) {
		currentNodeStruct := queue[head]
		head++
		currentElementName := currentNodeStruct.Name

		if expandedInThisReconstruction[currentElementName] { continue }
		expandedInThisReconstruction[currentElementName] = true

		if isSegmentToMeetingNode {
			if currentElementName == stopAtNode && currentElementName != startNode {
				currentNodeStruct.Children = nil 
				continue
			}
		} else { 
			if scraper.IsBaseElement(currentElementName) {
				currentNodeStruct.Children = nil 
				continue
			}
		}
		if scraper.IsBaseElement(currentElementName) { 
		    currentNodeStruct.Children = nil
			continue
		}

		var combinationToUse scraper.Combination
		foundCombinationInParents := false
		if comb, ok := parents[currentElementName]; ok {
			combinationToUse = comb
			foundCombinationInParents = true
		}
		if !foundCombinationInParents {
			if recipesForNode, recipeExists := (*recipeData)[currentElementName]; recipeExists && len(recipesForNode) > 0 {
				combinationToUse = recipesForNode[0] 
			} else {
				currentNodeStruct.Children = nil
				continue
			}
		}
		
		ing1Name := combinationToUse.First()
		ing2Name := combinationToUse.Second()

		if ing1Name == "" && ing2Name == "" && !scraper.IsBaseElement(currentElementName) {
			currentNodeStruct.Children = nil 
			continue
		}
		
		if ing1Name != "" || ing2Name != "" { 
			plusNode := scraper.TreeNode{Name: "+"}
			if ing1Name != "" {
				childNode1 := scraper.TreeNode{Name: ing1Name}
				plusNode.Children = append(plusNode.Children, childNode1)
				if !expandedInThisReconstruction[ing1Name] { 
					queue = append(queue, &plusNode.Children[len(plusNode.Children)-1])
				}
			}
			if ing2Name != "" {
				childNode2 := scraper.TreeNode{Name: ing2Name}
				plusNode.Children = append(plusNode.Children, childNode2)
				if !expandedInThisReconstruction[ing2Name] {
					queue = append(queue, &plusNode.Children[len(plusNode.Children)-1])
				}
			}
			currentNodeStruct.Children = append(currentNodeStruct.Children, plusNode)
		} else {
			currentNodeStruct.Children = nil 
		}
	}
	return root
}

func BidirectionalSearch(recipe *scraper.Recipe, tierMap *scraper.Tier, start string, meetingNodeChoiceIndex int) scraper.TreeNode {
	if recipe == nil {
		return scraper.TreeNode{Name: "Error: Recipe data is nil."}
	}
	isTargetBase := scraper.IsBaseElement(start)
	if isTargetBase {
		return scraper.TreeNode{Name: start}
	}
	if _, targetExistsInRecipe := (*recipe)[start]; !targetExistsInRecipe {
		return scraper.TreeNode{Name: fmt.Sprintf("Target element '%s' not found in recipes.", start)}
	}
	if meetingNodeChoiceIndex <= 0 {
		meetingNodeChoiceIndex = 1
	}

	// Inisialisasi
	visitedForward := make(map[string]bool)
	parentsForward := make(map[string]scraper.Combination)
	qForward := make([]string, 0)
	baseElementsSlice := getBaseElements(tierMap)
	if len(baseElementsSlice) == 0 {
		return scraper.TreeNode{Name: "Error: No base elements found."}
	}
	for _, el := range baseElementsSlice {
		if !visitedForward[el] {
			qForward = append(qForward, el)
			visitedForward[el] = true
		}
	}

	visitedBackward := make(map[string]bool)
	parentsBackward := make(map[string]scraper.Combination)
	qBackward := make([]string, 0)
	qBackward = append(qBackward, start)
	visitedBackward[start] = true

	collectedMeetingNodesMap := make(map[string]bool)
	var orderedMeetingNodes []string

	currentForwardLayer := qForward
	currentBackwardLayer := qBackward
	qForward = nil
	qBackward = nil

	maxSearchDepth := 20
	currentDepth := 0

	for (len(currentForwardLayer) > 0 || len(currentBackwardLayer) > 0) && currentDepth < maxSearchDepth {
		currentDepth++
		// Forward step
		if len(currentForwardLayer) > 0 {
			nextForwardLayer := []string{}
			for _, elementFwd := range currentForwardLayer {
				if visitedBackward[elementFwd] {
					if !collectedMeetingNodesMap[elementFwd] {
						collectedMeetingNodesMap[elementFwd] = true
						orderedMeetingNodes = append(orderedMeetingNodes, elementFwd)
					}
				}

				for product, combinations := range *recipe {
					if visitedForward[product] {
						continue
					}
					for _, comb := range combinations {
						ing1, ing2 := comb.First(), comb.Second()
						if visitedForward[ing1] && visitedForward[ing2] {
							if !parentsForwardKnown(parentsForward, product) {
								parentsForward[product] = comb
							}
							if !contains(nextForwardLayer, product) && !visitedForward[product] {
								nextForwardLayer = append(nextForwardLayer, product)
							}
						}
					}
				}
			}
			for _, node := range nextForwardLayer {
				visitedForward[node] = true
			}
			currentForwardLayer = nextForwardLayer
		}

		// Backward step
		if len(currentBackwardLayer) > 0 {
			nextBackwardLayer := []string{}
			for _, elementBwd := range currentBackwardLayer {
				if visitedForward[elementBwd] {
					if !collectedMeetingNodesMap[elementBwd] {
						collectedMeetingNodesMap[elementBwd] = true
						orderedMeetingNodes = append(orderedMeetingNodes, elementBwd)
					}
				}
				if scraper.IsBaseElement(elementBwd) {
					continue
				}
				combinations, found := (*recipe)[elementBwd]
				if found && len(combinations) > 0 {
					chosenCombination := combinations[0]
					if !parentsBackwardKnown(parentsBackward, elementBwd) {
						parentsBackward[elementBwd] = chosenCombination
					}
					ing1, ing2 := chosenCombination.First(), chosenCombination.Second()
					if ing1 != "" && !visitedBackward[ing1] {
						visitedBackward[ing1] = true
						if !contains(nextBackwardLayer, ing1) {
							nextBackwardLayer = append(nextBackwardLayer, ing1)
						}
					}
					if ing2 != "" && !visitedBackward[ing2] {
						visitedBackward[ing2] = true
						if !contains(nextBackwardLayer, ing2) {
							nextBackwardLayer = append(nextBackwardLayer, ing2)
						}
					}
				}
			}
			currentBackwardLayer = nextBackwardLayer
		}
		 if len(orderedMeetingNodes) >= meetingNodeChoiceIndex && len(currentForwardLayer) == 0 && len(currentBackwardLayer) == 0 {
            break
        }
	}
	
	sort.Strings(orderedMeetingNodes)

	var chosenMeetingNode string
	if len(orderedMeetingNodes) > 0 {
		if meetingNodeChoiceIndex > len(orderedMeetingNodes) {
			chosenMeetingNode = orderedMeetingNodes[len(orderedMeetingNodes)-1]
		} else {
			chosenMeetingNode = orderedMeetingNodes[meetingNodeChoiceIndex-1]
		}
	} else {
		return scraper.TreeNode{Name: fmt.Sprintf("No path found to '%s' (no meeting nodes found).", start)}
	}

	pathFromStartToMeetingNode := reconstructPathDeterministic(start, parentsBackward, recipe, chosenMeetingNode, true)
	
	var subTreeFromMeetingNodeToBases scraper.TreeNode
	if scraper.IsBaseElement(chosenMeetingNode) {
		subTreeFromMeetingNodeToBases = scraper.TreeNode{Name: chosenMeetingNode}
	} else {
		subTreeFromMeetingNodeToBases = reconstructPathDeterministic(chosenMeetingNode, parentsForward, recipe, "", false)
	}

	if start == chosenMeetingNode {
		pathFromStartToMeetingNode.Children = subTreeFromMeetingNodeToBases.Children
		return pathFromStartToMeetingNode
	}
	
	nodeToAttach := findNodeInTreeBFS(&pathFromStartToMeetingNode, chosenMeetingNode)
	if nodeToAttach != nil {
		nodeToAttach.Children = subTreeFromMeetingNodeToBases.Children
	} else {
		return scraper.TreeNode{Name: fmt.Sprintf("Error: Path reconstruction failed for '%s' at meeting node '%s'.", start, chosenMeetingNode)}
	}

	return pathFromStartToMeetingNode
}

func parentsForwardKnown(parents map[string]scraper.Combination, product string) bool {
	_, ok := parents[product]
	return ok
}
func parentsBackwardKnown(parents map[string]scraper.Combination, element string) bool {
	_, ok := parents[element]
	return ok
}
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}