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
	return s == "Air" || s == "Earth" || s == "Fire" || s == "Water" || s == "Time"
}

// Recipe is a map where the key is a string and the value is a slice of Combination
// representing the combinations of elements that can be made from the key element.
type Recipe map[string][]Combination

type Tier map[string]int

// TreeNode is the struct for frontend
// representation of the recipe tree.
type TreeNode struct {
	Name     string     `json:"name"`
	Children []TreeNode `json:"children"`
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

func (recipe* Recipe) SortRecipeChildren(tier* Tier){
	for key, combinations := range *recipe {
		value := []int{}
		for _, combination := range combinations {
			first := combination.First()
			second := combination.Second()
			value = append(value, max((*tier)[first], (*tier)[second]))
		}
		// Sort combinations based on the corresponding value
		sort.Slice(combinations, func(i, j int) bool {
			return value[i] < value[j]
		})
		(*recipe)[key] = combinations
	}
}