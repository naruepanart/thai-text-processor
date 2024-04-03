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
	// Open and parse blacklist JSON
	blacklist, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("failed to open blacklist.json: %v", err))
	}
	defer blacklist.Close()

	var replacements map[string]string
	if err := json.NewDecoder(blacklist).Decode(&replacements); err != nil {
		panic(fmt.Errorf("failed to decode blacklist.json: %v", err))
	}

	// Compile regex patterns
	patterns := make([]*regexp.Regexp, 0, len(replacements))
	values := make([]string, 0, len(replacements))
	for pattern, value := range replacements {
		patterns = append(patterns, regexp.MustCompile(pattern))
		values = append(values, value)
	}

	// Process text files
	txtFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	for _, txtFile := range txtFiles {
		// Open file
		file, err := os.Open(txtFile)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to open file %s: %v", txtFile, err))
			continue
		}
		defer file.Close()

		var updated []byte
		buf := make([]byte, 1024)

		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(fmt.Errorf("failed to read file %s: %v", txtFile, err))
				break
			}

			// Process and update text
			text := buf[:n]
			for i, pattern := range patterns {
				text = pattern.ReplaceAll(text, []byte(values[i]))
			}

			updated = append(updated, text...)
		}

		// Write updated content to file
		err = os.WriteFile(txtFile, updated, os.ModePerm)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to write file %s: %v", txtFile, err))
			continue
		}

		fmt.Printf("Text in %s updated.\n", txtFile)
	}
}