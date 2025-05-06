package scraper

import (
	"encoding/json"
	"log"
	"os"
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

func JsonToMap(filename string) *Recipe {
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

// TreeNode is the struct for frontend
// representation of the recipe tree.
type TreeNode struct {
	Name     string     `json:"name"`
	Children []TreeNode `json:"children"`
}
