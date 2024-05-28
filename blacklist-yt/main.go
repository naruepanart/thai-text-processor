package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	// Load the blacklist
	blacklistFile, err := os.Open("blacklist.json")
	if err != nil {
		fmt.Printf("Error loading blacklist: %v\n", err)
		return
	}
	defer blacklistFile.Close()

	var blacklist map[string]string
	if err := json.NewDecoder(blacklistFile).Decode(&blacklist); err != nil {
		fmt.Printf("Error decoding blacklist: %v\n", err)
		return
	}

	// Compile the blacklist regex patterns and replacements
	var regexPatterns []*regexp.Regexp
	var replacements [][]byte

	for pattern, replacement := range blacklist {
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Printf("Error compiling regex: %v\n", err)
			return
		}
		regexPatterns = append(regexPatterns, re)
		replacements = append(replacements, []byte(replacement))
	}

	// Find all text files
	textFiles, err := filepath.Glob("*.txt")
	if err != nil {
		fmt.Printf("Error finding text files: %v\n", err)
		return
	}

	// Process each text file
	for _, filename := range textFiles {
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			continue
		}

		updatedContent := content
		for i, pattern := range regexPatterns {
			updatedContent = pattern.ReplaceAll(updatedContent, replacements[i])
		}

		// Convert B.E. years to A.D. if they start with "25"
		re := regexp.MustCompile(`(\d{4})`)
		updatedContent = re.ReplaceAllFunc(updatedContent, func(match []byte) []byte {
			yearBE, _ := strconv.Atoi(string(match))
			if yearBE >= 2400 {
				yearAD := yearBE - 543
				return []byte(strconv.Itoa(yearAD))
			}
			return match
		})

		if err := os.WriteFile(filename, updatedContent, os.ModePerm); err != nil {
			fmt.Printf("Error writing to file %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Text in %s updated.\n", filename)
	}
}
