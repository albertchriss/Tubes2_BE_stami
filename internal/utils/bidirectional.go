package utils

import (
	"fmt"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
)

func findNodeInTree(root *scraper.TreeNode, name string) *scraper.TreeNode {
	if root == nil {
		return nil
	}
	if root.Name == name {
		return root
	}

	queue := []*scraper.TreeNode{}
	for i := range root.Children {
		queue = append(queue, &root.Children[i])
	}

	head := 0
	for head < len(queue) {
		currNode := queue[head]
		head++

		if currNode.Name == name {
			return currNode
		}

		for i := range currNode.Children {
			queue = append(queue, &currNode.Children[i])
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

func BidirectionalSearch(recipe *scraper.Recipe, tierMap *scraper.Tier, start string) scraper.TreeNode {
	if recipe == nil {
		return scraper.TreeNode{Name: "Error: Recipe data is nil."}
	}

	_, targetExistsInRecipe := (*recipe)[start]
	isTargetBase := scraper.IsBaseElement(start)

	if !targetExistsInRecipe && !isTargetBase {
		return scraper.TreeNode{Name: fmt.Sprintf("Target element '%s' not found in recipes and is not a base element.", start)}
	}

	if isTargetBase {
		return scraper.TreeNode{Name: start}
	}

	// Inisialisasi Forward Search
	visitedForward := make(map[string]bool)
	parentsForward := make(map[string]scraper.Combination)
	currentForwardLayer := make([]string, 0)

	baseElementsSlice := getBaseElements(tierMap)
	if len(baseElementsSlice) == 0 {
		return scraper.TreeNode{Name: "Error: No base elements found for forward search."}
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
								// fmt.Printf("Meeting node found by forward search: %s\n", meetingNode)
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

		if len(queueBackward) > 0 && meetingNode == "" {
			currentElementBackward := queueBackward[0]
			queueBackward = queueBackward[1:]

			if scraper.IsBaseElement(currentElementBackward) {
				if visitedForward[currentElementBackward] {
					meetingNode = currentElementBackward
					// fmt.Printf("Meeting node found by backward search (is base element): %s\n", meetingNode)
					break 
				}
				continue
			}

			combinations, found := (*recipe)[currentElementBackward]
			if !found || len(combinations) == 0 {
				continue
			}

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


	if meetingNode == "" {
		return scraper.TreeNode{Name: fmt.Sprintf("No path found to '%s'", start)}
	}
	
	pathFromStartToMeetingNode := reconstructPath(start, parentsBackward, recipe, meetingNode, false)

	var subTreeFromMeetingNodeToBases scraper.TreeNode
	if scraper.IsBaseElement(meetingNode) {
		subTreeFromMeetingNodeToBases = scraper.TreeNode{}
	} else {
		subTreeFromMeetingNodeToBases = reconstructPath(meetingNode, parentsForward, recipe, "", true)
	}

	meetingNodeInPath := findNodeInTree(&pathFromStartToMeetingNode, meetingNode)

	if meetingNodeInPath != nil {
		meetingNodeInPath.Children = subTreeFromMeetingNodeToBases.Children
	} else if start == meetingNode { 
		pathFromStartToMeetingNode.Children = subTreeFromMeetingNodeToBases.Children
	}

	return pathFromStartToMeetingNode
}

func reconstructPath(startNode string, parents map[string]scraper.Combination, recipeData *scraper.Recipe, stopAtNode string, isForwardStyle bool) scraper.TreeNode {
    root := scraper.TreeNode{Name: startNode}
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

        if stopAtNode != "" && currNode.Name == stopAtNode && currNode.Name != startNode {
            currNode.Children = nil
            continue
        }
        if scraper.IsBaseElement(currNode.Name) {
            currNode.Children = nil
            continue
        }

        var combinationToUse scraper.Combination
        foundCombination := false

        if comb, ok := parents[currNode.Name]; ok {
            combinationToUse = comb
            foundCombination = true
        } else if !isForwardStyle {
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

            if ing1Name != "" {
                childNode1 := scraper.TreeNode{Name: ing1Name}
                childrenForPlusNode = append(childrenForPlusNode, childNode1)
            }
            if ing2Name != "" {
                childNode2 := scraper.TreeNode{Name: ing2Name}
                childrenForPlusNode = append(childrenForPlusNode, childNode2)
            }

            if len(childrenForPlusNode) > 0 {
                plusNode.Children = childrenForPlusNode
                currNode.Children = append(currNode.Children, plusNode)

                for i := range plusNode.Children {
                    childToExpand := &plusNode.Children[i]
                    if !visitedExpansion[childToExpand.Name] {
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