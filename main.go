package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// TextProc represents a text processing utility.
type TextProc struct{ Replacements map[string]string }

// Process replaces words in the given text based on predefined replacements.
func (p TextProc) Process(text string) string {
	for word, replacement := range p.Replacements {
		text = regexp.MustCompile(word).ReplaceAllString(text, replacement)
	}
	return text
}

// FileMgr manages file-related operations using a TextProc.
type FileMgr struct{ Processor TextProc }

// ProcessFile reads, processes, and writes the updated text to a file.
func (m FileMgr) ProcessFile(filename string) error {
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

// InitProc initializes a TextProc with replacements from a JSON file.
func InitProc() TextProc {
	file, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("error opening blacklist.json: %v", err))
	}
	defer file.Close()

	var replacements map[string]string
	if err := json.NewDecoder(file).Decode(&replacements); err != nil {
		panic(fmt.Errorf("error decoding blacklist.json: %v", err))
	}

	return TextProc{Replacements: replacements}
}

func main() {
	textFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	manager := FileMgr{Processor: InitProc()}
	for _, filename := range textFiles {
		if err := manager.ProcessFile(filename); err != nil {
			fmt.Println(err)
		}
	}
}