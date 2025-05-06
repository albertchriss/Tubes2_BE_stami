package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Scraper(recipeFilename string, tierFileName string) error {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

	outputDir := filepath.Dir(recipeFilename)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		log.Printf("Output directory %s does not exist, creating...", outputDir)
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
		}
		log.Printf("Output directory %s created successfully.", outputDir)
	} else if err != nil {
		return fmt.Errorf("error checking output directory %s: %w", outputDir, err)
	}
    log.Println("Starting data scraping...")

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, http.StatusText(res.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	elementsData := make(map[string][][]string)
	elementsTier := make(map[string]int)

	// class "list-table col-list icon-hover" tbody tr
	doc.Find("table.list-table.col-list.icon-hover tbody").Each(func(num int, table *goquery.Selection) {
		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			// skip header
			if i == 0 {
				return
			}
	
			elementName := ""
			firstTD := row.Find("td").First()
			if firstTD.Length() > 0 {
				elementName = strings.TrimSpace(firstTD.Find("a").Text())
			}
	
			var recipeCombinations [][]string
			secondTD := row.Find("td").Eq(1)
			if secondTD.Length() > 0 {
				recipesList := secondTD.Find("ul li")
	
				recipesList.Each(func(j int, recipeItem *goquery.Selection) {
					var ingredients []string
					recipeItem.Find("a").Each(func(k int, ingredientLink *goquery.Selection) {
						ingredientName := strings.TrimSpace(ingredientLink.Text())
						if ingredientName != "" {
							ingredients = append(ingredients, ingredientName)
						}
					})
					if len(ingredients) > 0 {
						recipeCombinations = append(recipeCombinations, ingredients)
					}
				})
			}
	
			if elementName != "" {
				if len(recipeCombinations) > 0 {
					elementsData[elementName] = recipeCombinations
				} else {
					elementsData[elementName] = [][]string{}
				}

				if num == 0 {
					elementsTier[elementName] = 0
				} else {
					elementsTier[elementName] = num-1
				}
			}
		})

	})

	// save as JSON
	jsonData, err := json.MarshalIndent(elementsData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	
	file, err := os.Create(recipeFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", recipeFilename, err)
	}
	defer file.Close()
	
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file %s: %w", recipeFilename, err)
	}
	
	jsonTierData, err := json.MarshalIndent(elementsTier, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	tierFile, err := os.Create(tierFileName)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", tierFileName, err)
	}
	defer tierFile.Close()
	_, err = tierFile.Write(jsonTierData)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file %s: %w", tierFileName, err)
	}

	// fmt.Println(string(jsonData))
	fmt.Printf("Data saved to %s\n", recipeFilename)
	fmt.Printf("Tier data saved to %s\n", tierFileName)

	return nil
}