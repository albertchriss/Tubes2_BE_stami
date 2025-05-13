package scraper

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

type Combination [2]string

func (c Combination) First() string {
	if len(c) > 0 {
		return c[0]
	}
	return ""
}
func (c Combination) Second() string {
	if len(c) > 1 {
		return c[1]
	}
	return ""
}

func IsBaseElement(s string) bool {
	return s == "Air" || s == "Earth" || s == "Fire" || s == "Water"
}

// Recipe is a map where the key is a string and the value is a slice of Combination
// representing the combinations of elements that can be made from the key element.
type Recipe map[string][]Combination

type ElementInfo struct {
	Tier     int    `json:"tier"`
	ImageSrc string `json:"imageSrc"`
}

type Tier map[string]ElementInfo

// TreeNode is the struct for frontend
// representation of the recipe tree.
type TreeNode struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	ImageSrc string     `json:"imageSrc"`
	Children []TreeNode `json:"children"`
}

type SearchResult struct {
	Tree      TreeNode
	NodeCount int
	TimeTaken int64
}

func JsonToRecipe(filename string) *Recipe {
	var result Recipe
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	err = json.Unmarshal(file, &result)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	log.Println("JSON file successfully converted to map")
	return &result
}

func JsonToTier(filename string) *Tier {
	var result Tier
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	err = json.Unmarshal(file, &result)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}
	log.Println("JSON file successfully converted to map")
	return &result
}

type ElementTier struct {
	Name string
	Tier int
}

func (recipe *Recipe) SortRecipeChildren(tier *Tier) {
	sortedTier := []ElementTier{}
	for name, info := range *tier {
		num := info.Tier
		sortedTier = append(sortedTier, ElementTier{Name: name, Tier: num})
	}

	sort.Slice(sortedTier, func(i, j int) bool {
		return sortedTier[i].Tier < sortedTier[j].Tier
	})

	for _, element := range sortedTier {
		key := element.Name
		combinations := (*recipe)[key]
		newCombs := []Combination{}
		value := [][]int{}
		for _, combination := range combinations {
			_, foundFirst := (*tier)[combination.First()]
			if !foundFirst {
				continue
			}
			_, foundSec := (*tier)[combination.Second()]
			if !foundSec {
				continue
			}

			first := combination.First()
			second := combination.Second()

			if (*tier)[first].Tier >= (*tier)[key].Tier || ((*tier)[second].Tier >= (*tier)[key].Tier) {
				continue
			}

			maks := max((*tier)[first].Tier, (*tier)[second].Tier)
			mini := min((*tier)[first].Tier, (*tier)[second].Tier)
			value = append(value, []int{maks, mini})
			newCombs = append(newCombs, combination)
		}
		// Sort combinations based on the corresponding value
		if len(newCombs) > 0 {
			indices := make([]int, len(newCombs))
			for k := range indices {
				indices[k] = k
			}
			sort.Slice(indices, func(i, j int) bool {
				if value[indices[i]][0] == value[indices[j]][0] {
					return value[indices[i]][1] < value[indices[j]][1]
				}

				return value[indices[i]][0] < value[indices[j]][0]
			})
			sortedCombs := make([]Combination, len(newCombs))
			for i, index := range indices {
				sortedCombs[i] = newCombs[index]
			}
			newCombs = sortedCombs
			(*recipe)[key] = newCombs
		} else {
			if !IsBaseElement(key) {
				delete(*recipe, key)
				delete(*tier, key)
			}
		}
	}

}

// helper function
func (currNode *TreeNode) CountNumRecipe(parent *map[int]*TreeNode, numPath *map[int]int) int {
	num := 0

	if currNode.Name == "+" {
		num = (*numPath)[currNode.Children[0].Id] * (*numPath)[currNode.Children[1].Id]
	} else {
		for _, child := range currNode.Children {
			num += (*numPath)[child.Id]
		}
	}
	(*numPath)[currNode.Id] = num

	if currNode.Id == 0 {
		return num
	}

	par := (*parent)[currNode.Id]
	return par.CountNumRecipe(parent, numPath)
}

// node has to be a "+" node
func (node *TreeNode) InitParAndNum(parNode *TreeNode, parent *map[int]*TreeNode, numPath *map[int]int) {
	if node.Name != "+" {
		return
	}

	(*parent)[node.Id] = parNode
	(*parent)[node.Children[0].Id] = node
	(*parent)[node.Children[1].Id] = node
	(*numPath)[node.Id] = 1
	(*numPath)[node.Children[0].Id] = 1
	(*numPath)[node.Children[1].Id] = 1
}
