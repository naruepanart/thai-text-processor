package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Replacements map[string]string

func processFile(filename string, repl Replacements) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open file %s: %v", filename, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat file %s: %v", filename, err)
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return fmt.Errorf("read file %s: %v", filename, err)
	}

	updatedText := processText(string(data), repl)
	if err := os.WriteFile(filename, []byte(updatedText), 0644); err != nil {
		return fmt.Errorf("write file %s: %v", filename, err)
	}

	fmt.Printf("Text in %s updated.\n", filename)
	return nil
}

func processText(text string, repl Replacements) string {
	for pattern, replacement := range repl {
		text = regexp.MustCompile(pattern).ReplaceAllString(text, replacement)
	}
	return text
}

func loadReplacements(filename string) Replacements {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("open %s: %v", filename, err))
	}
	defer file.Close()

	var repl Replacements
	if err := json.NewDecoder(file).Decode(&repl); err != nil {
		panic(fmt.Errorf("decode %s: %v", filename, err))
	}

	return repl
}

func main() {
	txtFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	repl := loadReplacements("blacklist.json")
	repl[`(\d+)-(\d+)`] = `${1}ถึง${2}`

	for _, filename := range txtFiles {
		if err := processFile(filename, repl); err != nil {
			fmt.Println(err)
		}
	}
}
