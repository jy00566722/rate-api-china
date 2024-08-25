package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const localesDir = "/Users/pengfeng/vscode/internation-rate/internation-wxt-vue-rate/assets/locales1"

func main() {
	// Load .env environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// test()

	// Walk through the directory
	filepath.Walk(localesDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Check for JSON files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			for _, langItem := range languageList {
				// Match JSON file with language code
				if strings.EqualFold(fmt.Sprintf("%s.json", langItem.Code), info.Name()) {
					// Read the JSON file
					file, err := os.ReadFile(path)
					if err != nil {
						log.Fatalf("failed to read file: %v", err)
					}

					// Parse the JSON content
					var content map[string]interface{}
					json.Unmarshal(file, &content)

					// Use the translate function to get the translated content
					translated, err := translate("Setting theme", "en", langItem.Bd)
					if err != nil {
						log.Printf("failed to translate for language '%s': %v", langItem.Code, err)
						continue
					}

					// Update the JSON content
					content["settingTheme"] = translated

					// Marshal the content back to JSON
					updatedJSON, err := json.MarshalIndent(content, "", "  ")
					if err != nil {
						log.Printf("failed to marshal JSON: %v", err)
						continue
					}

					// Write the updated JSON back into the file
					err = os.WriteFile(path, updatedJSON, 0644)
					if err != nil {
						log.Printf("failed to write updated JSON: %v", err)
						continue
					}

					fmt.Printf("Successfully updated file: %s\n", info.Name())
					time.Sleep(2 * time.Second)
					break
				}
			}
		}

		return nil
	})
}
