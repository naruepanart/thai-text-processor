package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// TextProcessor represents a text processing utility.
type TextProcessor struct{ Replacements map[string]string }

// Process replaces words in the given text based on predefined replacements.
func (p TextProcessor) Process(text string) string {
	for word, replacement := range p.Replacements {
		text = regexp.MustCompile(word).ReplaceAllString(text, replacement)
	}
	return text
}

// FileManager manages file-related operations using a TextProcessor.
type FileManager struct{ Processor TextProcessor }

// ProcessFile reads, processes, and writes the updated text to a file.
func (m FileManager) ProcessFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	updatedText := m.Processor.Process(string(data))
	if err := os.WriteFile(filename, []byte(updatedText), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %v", filename, err)
	}

	fmt.Printf("Text in %s has been updated.\n", filename)
	return nil
}

// initProcessor initializes a TextProcessor with replacements from a JSON file.
func initProcessor() TextProcessor {
	file, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("error opening blacklist.json: %v", err))
	}
	defer file.Close()

	var replacements map[string]string
	if err := json.NewDecoder(file).Decode(&replacements); err != nil {
		panic(fmt.Errorf("error decoding blacklist.json: %v", err))
	}

	return TextProcessor{Replacements: replacements}
}

func main() {
	txtFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	manager := FileManager{Processor: initProcessor()}
	for _, fileName := range txtFiles {
		if err := manager.ProcessFile(fileName); err != nil {
			fmt.Println(err)
		}
	}
}