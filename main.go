package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// TextProcessor represents a text processing object.
type TextProcessor struct {
	Replacements map[string]string
}

// Process applies replacements to the input text.
func (tp TextProcessor) Process(text string) string {
	for pattern, replacement := range tp.Replacements {
		text = regexp.MustCompile(pattern).ReplaceAllString(text, replacement)
	}
	return text
}

// FileManager manages file operations with a TextProcessor.
type FileManager struct {
	Processor TextProcessor
}

// ProcessFile reads, processes, and writes the updated text to a file.
func (fm FileManager) ProcessFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	updatedText := fm.Processor.Process(string(data))
	if err := os.WriteFile(filename, []byte(updatedText), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %v", filename, err)
	}

	fmt.Printf("Text in %s updated.\n", filename)
	return nil
}

// NewProcessor initializes a TextProcessor with replacements from "blacklist.json".
func NewProcessor() TextProcessor {
	file, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("error opening blacklist.json: %v", err))
	}
	defer file.Close()

	var replacements map[string]string
	if err := json.NewDecoder(file).Decode(&replacements); err != nil {
		panic(fmt.Errorf("error decoding blacklist.json: %v", err))
	}

	// Additional replacements
	replacements[`(\d+)-(\d+)`] = `${1}ถึง${2}`
	replacements[`(\p{L}+)\s+\(`] = `${1} หรือ `
	replacements[`(\s*\))`] = ``

	return TextProcessor{Replacements: replacements}
}

func main() {
	txtFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	fileManager := FileManager{Processor: NewProcessor()}
	for _, filename := range txtFiles {
		if err := fileManager.ProcessFile(filename); err != nil {
			fmt.Println(err)
		}
	}
}