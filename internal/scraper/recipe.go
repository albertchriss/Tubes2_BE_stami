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

func (recipe *Recipe) SortRecipeChildren(tier *Tier) {
	for key, combinations := range *recipe {
		newCombs := []Combination{}
		value := [][]int{}
		for _, combination := range combinations {
			_, foundFirst := (*tier)[combination.First()]
			if !foundFirst {
				log.Printf("Element %s not found in tier map\n", combination.First())
				continue
			}
			_, foundSec := (*tier)[combination.Second()]
			if !foundSec {
				log.Printf("Element %s not found in tier map\n", combination.Second())
				continue
			}

			first := combination.First()
			second := combination.Second()
			maks := max((*tier)[first], (*tier)[second])
			mini := min((*tier)[first], (*tier)[second])
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
		}
		(*recipe)[key] = newCombs
	}
}
