package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	var blacklistReplacements map[string]string
	if err := json.NewDecoder(blacklistFile).Decode(&blacklistReplacements); err != nil {
		panic(fmt.Errorf("failed to decode blacklist.json: %v", err))
	}

	// Precompile regular expressions
	compiledReplacements := make(map[*regexp.Regexp]string)
	for pattern, replacement := range blacklistReplacements {
		compiledReplacements[regexp.MustCompile(pattern)] = replacement
	}

	// Process all text files in the current directory
	txtFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	for _, txtFile := range txtFiles {
		// Read file content in chunks to minimize heap allocations
		file, err := os.Open(txtFile)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to open file %s: %v", txtFile, err))
			continue
		}
		defer file.Close()

		var updatedTextBuilder []byte
		buf := make([]byte, 1024) // Adjust the buffer size based on your file characteristics

		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(fmt.Errorf("failed to read file %s: %v", txtFile, err))
				break
			}

			// Process text and update file content
			updatedText := string(buf[:n])
			for pattern, replacement := range compiledReplacements {
				updatedText = pattern.ReplaceAllString(updatedText, replacement)
			}

			updatedTextBuilder = append(updatedTextBuilder, []byte(updatedText)...)
		}

		// Use stack memory to write updated content to file
		err = os.WriteFile(txtFile, updatedTextBuilder, os.ModePerm)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to write file %s: %v", txtFile, err))
			continue
		}

		fmt.Printf("Text in %s updated.\n", txtFile)
	}
}