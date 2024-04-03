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
	// Open and parse the blacklist JSON file
	f, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("open blacklist.json: %v", err))
	}
	defer f.Close()

	blacklist := make(map[string]string)
	if err := json.NewDecoder(f).Decode(&blacklist); err != nil {
		panic(fmt.Errorf("decode blacklist.json: %v", err))
	}

	// Compile regular expressions from the blacklist
	regexps := make([]*regexp.Regexp, 0, len(blacklist))
	replacements := make([][]byte, 0, len(blacklist))
	for pattern, replacement := range blacklist {
		re, err := regexp.Compile(pattern)
		if err != nil {
			panic(fmt.Errorf("compile regex %s: %v", pattern, err))
		}
		regexps = append(regexps, re)
		replacements = append(replacements, []byte(replacement))
	}

	// Find all text files in the current directory
	files, err := filepath.Glob("*.txt")
	if err != nil {
		panic(fmt.Errorf("glob: %v", err))
	}

	// Process each text file
	for _, filename := range files {
		file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println(fmt.Errorf("open file %s: %v", filename, err))
			continue
		}
		defer file.Close()

		// Process file content
		var updated []byte
		buf := make([]byte, 1024)

		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(fmt.Errorf("read file %s: %v", filename, err))
				break
			}

			// Apply regular expression replacements
			text := buf[:n]
			for i, re := range regexps {
				text = re.ReplaceAll(text, replacements[i])
			}

			updated = append(updated, text...)
		}

		// Write updated content back to the file
		if err := os.Truncate(filename, 0); err != nil {
			fmt.Println(fmt.Errorf("truncate file %s: %v", filename, err))
			continue
		}
		if _, err := file.WriteAt(updated, 0); err != nil {
			fmt.Println(fmt.Errorf("write file %s: %v", filename, err))
			continue
		}

		fmt.Printf("Text in %s updated.\n", filename)
	}
}