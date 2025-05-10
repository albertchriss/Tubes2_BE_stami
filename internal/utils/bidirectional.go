package utils

import (
	"fmt"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

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

func BidirectionalSearch(recipe *scraper.Recipe, tierMap *scraper.Tier, start string) scraper.BidirectionalResult {
	defaultResult := scraper.BidirectionalResult{PathFound: false}

	if recipe == nil {
		defaultResult.ForwardTree.Name = "Error: Recipe data is nil."
		defaultResult.BackwardTree.Name = "Error: Recipe data is nil."
		return defaultResult
	}

	_, targetExistsInRecipe := (*recipe)[start]
	isTargetBase := scraper.IsBaseElement(start)

	if !targetExistsInRecipe && !isTargetBase {
		defaultResult.ForwardTree.Name = fmt.Sprintf("Target element '%s' not found in recipes and is not a base element.", start)
		return defaultResult
	}

	if isTargetBase {
		return scraper.BidirectionalResult{
			ForwardTree:     scraper.TreeNode{Name: start},
			BackwardTree:    scraper.TreeNode{Name: start},
			MeetingNodeName: start,
			PathFound:       true,
		}
	}

	// Inisialisasi Forward Search
	visitedForward := make(map[string]bool)
	parentsForward := make(map[string]scraper.Combination)
	currentForwardLayer := make([]string, 0)

	baseElementsSlice := getBaseElements(tierMap)
	if len(baseElementsSlice) == 0 {
		defaultResult.ForwardTree.Name = "Error: No base elements found for forward search."
		return defaultResult
	}
	for _, el := range baseElementsSlice {
		if !visitedForward[el] {
			currentForwardLayer = append(currentForwardLayer, el)
			visitedForward[el] = true
		}
	}

	// Inisialisasi Backward Search
	visitedBackward := make(map[string]bool)
	parentsBackward := make(map[string]scraper.Combination)
	queueBackward := make([]string, 0)

	queueBackward = append(queueBackward, start)
	visitedBackward[start] = true

	meetingNode := ""

	for (len(currentForwardLayer) > 0 || len(queueBackward) > 0) && meetingNode == "" {
		if len(currentForwardLayer) > 0 && meetingNode == "" {
			nextForwardLayer := []string{}
			for product, combinations := range *recipe {
				if visitedForward[product] {
					continue
				}
				for _, comb := range combinations {
					ing1 := comb.First()
					ing2 := comb.Second()

					if visitedForward[ing1] && visitedForward[ing2] {
						if !visitedForward[product] {
							parentsForward[product] = comb
							visitedForward[product] = true
							nextForwardLayer = append(nextForwardLayer, product)

							if visitedBackward[product] {
								meetingNode = product
								break
							}
						}
					}
				}
				if meetingNode != "" {
					break
				}
			}
			currentForwardLayer = nextForwardLayer
		}

		if meetingNode != "" {
			break
		}

		// Backward Search Step
		if len(queueBackward) > 0 && meetingNode == "" {
			currentElementBackward := queueBackward[0]
			queueBackward = queueBackward[1:]

			if scraper.IsBaseElement(currentElementBackward) {
				if visitedForward[currentElementBackward] {
					meetingNode = currentElementBackward
					break
				}
				continue
			}

			combinations, found := (*recipe)[currentElementBackward]
			if !found {
				continue
			}

			if len(combinations) > 0 {
				chosenCombination := combinations[0]
				parentsBackward[currentElementBackward] = chosenCombination

				ing1 := chosenCombination.First()
				ing2 := chosenCombination.Second()

				processIngredient := func(ing string) bool {
					if ing != "" && !visitedBackward[ing] {
						visitedBackward[ing] = true
						queueBackward = append(queueBackward, ing)
						if visitedForward[ing] {
							meetingNode = ing
							return true
						}
					}
					return false
				}

				if processIngredient(ing1) {
					break
				}
				if meetingNode != "" {
					break
				}
				if processIngredient(ing2) {
					break
				}
			}
		}
	}

	if meetingNode == "" {
		defaultResult.ForwardTree.Name = fmt.Sprintf("No path found between base elements and '%s'.", start)
		defaultResult.BackwardTree.Name = fmt.Sprintf("No path found to '%s' from base elements or vice versa.", start)
		return defaultResult
	}

	result := scraper.BidirectionalResult{
		MeetingNodeName: meetingNode,
		PathFound:       true,
	}

	baseElementsMap := make(map[string]bool)
	for _, el := range baseElementsSlice {
		baseElementsMap[el] = true
	}

	result.ForwardTree = reconstructPathToMeetingNode(meetingNode, parentsForward, true, recipe, meetingNode)
	result.BackwardTree = reconstructPathToMeetingNode(start, parentsBackward, false, recipe, meetingNode)

	return result
}

func reconstructPathToMeetingNode(startNodeForPath string, parents map[string]scraper.Combination, isForwardPath bool, recipeData *scraper.Recipe, globalMeetingNode string) scraper.TreeNode {
	root := scraper.TreeNode{Name: startNodeForPath}
	queue := []*scraper.TreeNode{&root}
	visitedExpansion := make(map[string]bool)

	head := 0
	for head < len(queue) {
		currNode := queue[head]
		head++

		if visitedExpansion[currNode.Name] {
			continue
		}
		visitedExpansion[currNode.Name] = true

		if isForwardPath {
			if scraper.IsBaseElement(currNode.Name) {
				currNode.Children = nil
				continue
			}
		} else {
			if currNode.Name == globalMeetingNode && currNode.Name != startNodeForPath {
				currNode.Children = nil
				continue
			}
			if currNode.Name != globalMeetingNode && scraper.IsBaseElement(currNode.Name) {
				currNode.Children = nil
				continue
			}
		}

		var combinationToUse scraper.Combination
		foundCombination := false

		if comb, ok := parents[currNode.Name]; ok {
			combinationToUse = comb
			foundCombination = true
		} else if !scraper.IsBaseElement(currNode.Name) {
			if recipesForNode, recipeExists := (*recipeData)[currNode.Name]; recipeExists && len(recipesForNode) > 0 {
				combinationToUse = recipesForNode[0]
				foundCombination = true
			}
		}

		if !foundCombination {
			currNode.Children = nil
			continue
		}

		ing1Name := combinationToUse.First()
		ing2Name := combinationToUse.Second()

		if ing1Name == "" && ing2Name == "" && !scraper.IsBaseElement(currNode.Name) {
			currNode.Children = nil
			continue
		}

		if ing1Name != "" || ing2Name != "" {
			plusNode := scraper.TreeNode{Name: "+"}
			var childrenForPlusNode []scraper.TreeNode

			// Mengisi childrenForPlusNode
			if ing1Name != "" {
				childNode := scraper.TreeNode{Name: ing1Name}
				childrenForPlusNode = append(childrenForPlusNode, childNode)
			}
			if ing2Name != "" {
				childNode := scraper.TreeNode{Name: ing2Name}
				childrenForPlusNode = append(childrenForPlusNode, childNode)
			}

			if len(childrenForPlusNode) > 0 {
				plusNode.Children = childrenForPlusNode
				currNode.Children = append(currNode.Children, plusNode)

				// Tambahkan anak-anak dari plusNode ke queue utama untuk ekspansi
				for i := range plusNode.Children {
					childToExpand := &plusNode.Children[i]

					expandThisChild := true
					if isForwardPath {
						if scraper.IsBaseElement(childToExpand.Name) {
							expandThisChild = false
						}
					} else {
						if childToExpand.Name == globalMeetingNode && childToExpand.Name != startNodeForPath {
							expandThisChild = false
						} else if childToExpand.Name != globalMeetingNode && scraper.IsBaseElement(childToExpand.Name) {
							expandThisChild = false
						}
					}

					if expandThisChild && !visitedExpansion[childToExpand.Name] {
						queue = append(queue, childToExpand)
					}
				}
			} else {
				currNode.Children = nil
			}
		} else {
			currNode.Children = nil
		}
	}
	return root
}