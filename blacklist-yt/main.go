package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	// Load replacements from blacklist.json
	blacklistFile, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("failed to open blacklist.json: %v", err))
	}
	defer blacklistFile.Close()

	// Parse JSON content from blacklist file
	var replacements map[string]string
	if err := json.NewDecoder(blacklistFile).Decode(&replacements); err != nil {
		panic(fmt.Errorf("failed to decode blacklist.json: %v", err))
	}

	// Add custom replacement for number ranges
	replacements[`(\d+)-(\d+)`] = `${1}ถึง${2}`

	// Process all text files in the current directory
	txtFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	for _, filename := range txtFiles {
		// Read file content into memory
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to read file %s: %v", filename, err))
			continue
		}

		// Process text and update file content
		updatedText := string(fileContent)
		for pattern, replacement := range replacements {
			updatedText = regexp.MustCompile(pattern).ReplaceAllString(updatedText, replacement)
		}

		// Use stack memory to write updated content to file
		err = os.WriteFile(filename, []byte(updatedText), os.ModePerm)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to write file %s: %v", filename, err))
			continue
		}

		fmt.Printf("Text in %s updated.\n", filename)
	}
}