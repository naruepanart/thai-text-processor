package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type TextProcessor struct {
	Repl map[string]string
}

func (tp TextProcessor) Process(txt string) string {
	for pat, repl := range tp.Repl {
		txt = regexp.MustCompile(pat).ReplaceAllString(txt, repl)
	}
	return txt
}

type FileManager struct {
	Proc TextProcessor
}

func (fm FileManager) ProcessFile(fn string) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", fn, err)
	}

	updatedText := fm.Proc.Process(string(data))
	if err := os.WriteFile(fn, []byte(updatedText), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %v", fn, err)
	}

	fmt.Printf("Text in %s updated.\n", fn)
	return nil
}

func NewProcessor() TextProcessor {
	file, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("error opening blacklist.json: %v", err))
	}
	defer file.Close()

	var repl map[string]string
	if err := json.NewDecoder(file).Decode(&repl); err != nil {
		panic(fmt.Errorf("error decoding blacklist.json: %v", err))
	}

	repl[`(\d+)-(\d+)`] = `${1}ถึง${2}`
	repl[`(\p{L}+)\s+\(`] = `${1} หรือ `
	repl[`(\s*\))`] = ``

	return TextProcessor{Repl: repl}
}

func main() {
	txtFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	fm := FileManager{Proc: NewProcessor()}
	for _, fn := range txtFiles {
		if err := fm.ProcessFile(fn); err != nil {
			fmt.Println(err)
		}
	}
}