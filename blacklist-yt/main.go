package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	err := processFilesWithBlacklist()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func processFilesWithBlacklist() error {
	blacklist, err := loadBlacklist("blacklist.json")
	if err != nil {
		return fmt.Errorf("loading blacklist: %v", err)
	}

	regexPatterns, replacements, err := compileBlacklist(blacklist)
	if err != nil {
		return fmt.Errorf("compiling regex: %v", err)
	}

	textFiles, err := filepath.Glob("*.txt")
	if err != nil {
		return fmt.Errorf("glob: %v", err)
	}

	for _, filename := range textFiles {
		if err := processFileWithBlacklist(filename, regexPatterns, replacements); err != nil {
			fmt.Printf("Error processing %s: %v\n", filename, err)
		}
	}

	return nil
}

func loadBlacklist(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var blacklist map[string]string
	if err := json.NewDecoder(file).Decode(&blacklist); err != nil {
		return nil, err
	}
	return blacklist, nil
}

func compileBlacklist(blacklist map[string]string) ([]*regexp.Regexp, [][]byte, error) {
	var regexPatterns []*regexp.Regexp
	var replacements [][]byte

	for pattern, replacement := range blacklist {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, nil, err
		}
		regexPatterns = append(regexPatterns, re)
		replacements = append(replacements, []byte(replacement))
	}

	return regexPatterns, replacements, nil
}

func processFileWithBlacklist(filename string, regexPatterns []*regexp.Regexp, replacements [][]byte) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	updatedContent := content
	for i, pattern := range regexPatterns {
		updatedContent = pattern.ReplaceAll(updatedContent, replacements[i])
	}

	if err := os.WriteFile(filename, updatedContent, os.ModePerm); err != nil {
		return err
	}

	fmt.Printf("Text in %s updated.\n", filename)
	return nil
}