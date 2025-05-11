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
		for el, tier := range *tierMap {
			if tier == 0 {
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

func reconstructPathDeterministic(
	startNodeName string,
	parents map[string]scraper.Combination,
	recipeData *scraper.Recipe,
	stopAtNodeName string,
	isSegmentToMeetingNode bool,
) scraper.TreeNode {

	root := scraper.TreeNode{Name: startNodeName}
	if !isSegmentToMeetingNode && scraper.IsBaseElement(startNodeName) {
		return root
	}
	if isSegmentToMeetingNode && startNodeName == stopAtNodeName {
	}


	queue := []*scraper.TreeNode{&root}
	expanded := make(map[string]bool)

	for head := 0; head < len(queue); head++ {
		currNode := queue[head]
		name := currNode.Name

		if expanded[name] {
			continue
		}
		expanded[name] = true

		if isSegmentToMeetingNode {
			if name == stopAtNodeName && name != startNodeName {
				continue
			}
		} else {
			if scraper.IsBaseElement(name) {
				continue
			}
		}

		comb, ok := parents[name]
		if !ok {
			recipes, exists := (*recipeData)[name]
			if !exists || len(recipes) == 0 {
				continue
			}
			comb = recipes[0]
		}

		ing1, ing2 := comb.First(), comb.Second()
		if ing1 == "" && ing2 == "" {
			continue
		}

		tempChildrenForPlusNode := []scraper.TreeNode{}
		if ing1 != "" {
			tempChildrenForPlusNode = append(tempChildrenForPlusNode, scraper.TreeNode{Name: ing1})
		}
		if ing2 != "" {
			tempChildrenForPlusNode = append(tempChildrenForPlusNode, scraper.TreeNode{Name: ing2})
		}

		if len(tempChildrenForPlusNode) > 0 {
			plusNode := scraper.TreeNode{Name: "+", Children: tempChildrenForPlusNode}
			currNode.Children = append(currNode.Children, plusNode)

			addedPlusNodeRef := &currNode.Children[len(currNode.Children)-1]

			for i := range addedPlusNodeRef.Children {
				queue = append(queue, &addedPlusNodeRef.Children[i])
			}
		}
	}

	return root
}

func BidirectionalSearch(
	recipe *scraper.Recipe,
	tierMap *scraper.Tier,
	start string,
	meetingNodeChoiceIndex int,
) scraper.TreeNode {
	if recipe == nil {
		return scraper.TreeNode{Name: "Error: Recipe data is nil."}
	}
	if scraper.IsBaseElement(start) {
		return scraper.TreeNode{Name: start}
	}
	if _, exists := (*recipe)[start]; !exists {
		return scraper.TreeNode{Name: fmt.Sprintf("Target element '%s' not found in recipes.", start)}
	}

	if meetingNodeChoiceIndex < 1 {
		meetingNodeChoiceIndex = 1
	}

	visitedFwd := make(map[string]bool)
	parentsFwd := make(map[string]scraper.Combination)
	baseEls := getBaseElements(tierMap)

	if len(baseEls) == 0 {
		return scraper.TreeNode{Name: "Error: No base elements found."}
	}

	currentFwd := make([]string, 0, len(baseEls))
	for _, el := range baseEls {
		if !visitedFwd[el] {
			visitedFwd[el] = true
			currentFwd = append(currentFwd, el)
		}
	}

	visitedBwd := make(map[string]bool)
	parentsBwd := make(map[string]scraper.Combination)
	currentBwd := []string{start}
	visitedBwd[start] = true

	collectedMeetings := make(map[string]bool)
	orderedMeetings := []string{}

	const maxDepth = 20

	for depth := 0; depth < maxDepth; depth++ {
		nextFwd := []string{}
		for _, nodeFromBase := range currentFwd {
			if visitedBwd[nodeFromBase] && !collectedMeetings[nodeFromBase] {
				collectedMeetings[nodeFromBase] = true
				orderedMeetings = append(orderedMeetings, nodeFromBase)
			}

			for product, combinations := range *recipe {
				if visitedFwd[product] {
					continue
				}
				for _, comb := range combinations {
					ing1, ing2 := comb.First(), comb.Second()
					if (ing1 == ing2 && visitedFwd[ing1]) || (ing1 != ing2 && visitedFwd[ing1] && visitedFwd[ing2]) {
						if _, ok := parentsFwd[product]; !ok {
							parentsFwd[product] = comb
						}
						foundInNext := false
						for _, item := range nextFwd {
							if item == product {
								foundInNext = true
								break
							}
						}
						if !foundInNext {
							nextFwd = append(nextFwd, product)
						}
						break
					}
				}
			}
		}
		for _, n := range nextFwd {
			visitedFwd[n] = true
		}
		currentFwd = nextFwd

		if len(currentFwd) == 0 && len(orderedMeetings) == 0 && depth > 0 {
		}


		nextBwd := []string{}
		for _, nodeToMake := range currentBwd {
			if visitedFwd[nodeToMake] && !collectedMeetings[nodeToMake] {
				collectedMeetings[nodeToMake] = true
				orderedMeetings = append(orderedMeetings, nodeToMake)
			}

			if scraper.IsBaseElement(nodeToMake) {
				continue
			}

			recipesForNode, exists := (*recipe)[nodeToMake]
			if !exists || len(recipesForNode) == 0 {
				continue
			}
			
			chosenComb := recipesForNode[0]
			if _, ok := parentsBwd[nodeToMake]; !ok {
					parentsBwd[nodeToMake] = chosenComb
			}


			for _, ing := range []string{chosenComb.First(), chosenComb.Second()} {
				if ing != "" && !visitedBwd[ing] {
					visitedBwd[ing] = true
					foundInNext := false
					for _, item := range nextBwd {
						if item == ing {
							foundInNext = true
							break
						}
					}
					if !foundInNext {
						nextBwd = append(nextBwd, ing)
					}
				}
			}
		}
		currentBwd = nextBwd

		if len(currentBwd) == 0 && len(orderedMeetings) == 0 && depth > 0 {
		}

		if len(orderedMeetings) > 0 {
			break
		}

		if len(currentFwd) == 0 && len(currentBwd) == 0 && depth > 0 {
		    break
		}
	}

	if len(orderedMeetings) == 0 {
		return scraper.TreeNode{Name: fmt.Sprintf("No path found to '%s' (no meeting nodes discovered).", start)}
	}

	sort.Strings(orderedMeetings)
	
	var meeting string
	meetingIdx := meetingNodeChoiceIndex - 1

	if meetingIdx >= 0 && meetingIdx < len(orderedMeetings) {
		meeting = orderedMeetings[meetingIdx]
	} else if len(orderedMeetings) > 0 { 
		meeting = orderedMeetings[len(orderedMeetings)-1] 
	} else {
		return scraper.TreeNode{Name: fmt.Sprintf("Critical error: No meeting nodes available after selection for '%s'.", start)}
	}


	segStartToMeeting := reconstructPathDeterministic(start, parentsBwd, recipe, meeting, true)

	var segMeetingToBase scraper.TreeNode
	if scraper.IsBaseElement(meeting) {
		segMeetingToBase = scraper.TreeNode{}
	} else {
		segMeetingToBase = reconstructPathDeterministic(meeting, parentsFwd, recipe, "", false)
	}

	if start == meeting {
		segStartToMeeting.Children = segMeetingToBase.Children
		return segStartToMeeting
	}
	
	nodeToAttachChildren := findNodeInTreeBFS(&segStartToMeeting, meeting)
	if nodeToAttachChildren == nil {
		return scraper.TreeNode{Name: fmt.Sprintf("Error: Reconstruction failed. Meeting node '%s' not found in the first segment tree for target '%s'.", meeting, start)}
	}

	nodeToAttachChildren.Children = segMeetingToBase.Children
	
	return segStartToMeeting
}